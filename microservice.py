from flask import Flask, jsonify, request
import hvac
import os
import time

app = Flask(__name__)

# Retrieve Vault token from environment variable
vault_token = os.getenv('VAULT_TOKEN')
if not vault_token:
    raise ValueError("VAULT_TOKEN environment variable is not set")

# Vault configuration
vault_url = os.getenv('VAULT_ADDR', 'https://vault.primary-vault.svc.cluster.local:8200')
ca_cert_path = '/usr/local/share/ca-certificates/lab-root-ca.crt'
mount_point = 'hostnaming'

client = hvac.Client(url=vault_url, token=vault_token, verify=ca_cert_path)

def get_and_increment_counter(application, role, environment, max_retries=20):
    """
    Atomically reads and increments the counter using CAS.
    Returns the counter value for the current request.
    """
    counter_path = f"applications/{application}/{role}/{environment}/counter"
    retries = 0

    while retries < max_retries:
        try:
            counter_data = client.secrets.kv.v2.read_secret_version(
                path=counter_path, mount_point=mount_point
            )
            current_counter = counter_data["data"]["data"]["counter"]
            version = counter_data["data"]["metadata"]["version"]
        except hvac.exceptions.InvalidPath:
            # If the counter doesn't exist, create it with counter=2 so that the first value is 1.
            try:
                client.secrets.kv.v2.create_or_update_secret(
                    path=counter_path,
                    secret={"counter": 2},
                    mount_point=mount_point,
                    cas=0  # cas=0 means “only create if it doesn't exist”
                )
                return 1
            except Exception as e:
                retries += 1
                time.sleep(0.2)
                continue

        new_counter = current_counter + 1

        try:
            # Use CAS to update the counter atomically.
            client.secrets.kv.v2.create_or_update_secret(
                path=counter_path,
                secret={"counter": new_counter},
                mount_point=mount_point,
                cas=version
            )
            return current_counter  # Return the counter value before incrementing.
        except Exception as e:
            retries += 1
            time.sleep(0.2)
            continue

    raise Exception("Max retries reached in CAS update")

@app.route('/generate-name', methods=['POST'])
def generate_custom_name():
    try:
        data = request.json

        application = data.get("application")
        role = data.get("role")
        environment = data.get("environment")

        if not all([application, role, environment]):
            return jsonify({
                "error": "Missing one or more required fields: application, role, environment"
            }), 400

        # Atomically read and increment the global counter.
        counter = get_and_increment_counter(application, role, environment)
        hostname = f"{application}-{role}-{environment}-vm{counter}"

        return jsonify({"hostname": hostname}), 200

    except Exception as e:
        app.logger.error("Error generating hostname: %s", e)
        return jsonify({"error": "Internal server error. Please try again."}), 500

if __name__ == '__main__':
    # Run in non-debug mode to avoid abrupt termination on error.
    app.run(host='0.0.0.0', port=5000)
