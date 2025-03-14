# terraform-vault-hostnaming-poc

**WARNING:**  
This repository contains work-in-progress, proof-of-concept (PoC) code. It is not production-ready. Use it for testing and experimentation only.

---

## Overview

This repo demonstrates a PoC for dynamic host naming using Terraform and Vault. The solution consists of:

- **Vault Microservice:**  
  A Flask-based microservice that leverages Vault’s KV V2 secrets engine with Check-And-Set (CAS) for atomic counter updates.  
  - Generates unique hostnames based on parameters like _application_, _role_, and _environment_.  
  - Uses CAS to ensure global uniqueness even under concurrent updates.

- **Custom Terraform Provider:**  
  A custom Terraform provider that calls the microservice to generate hostnames.  
  - Simplifies integration with Terraform by exposing a resource that returns a generated hostname.  
  - Demonstrates two approaches:
    - **Example 1:** Using a `null_resource` with local-exec and an external data source to capture the output.
    - **Example 2:** Using the custom provider directly to provision naming resources.

---

## Components

### Vault Microservice

- **Language & Framework:** Python, Flask  
- **Key Features:**
  - Reads and increments a counter stored in Vault.
  - Uses CAS (Check-And-Set) for atomic updates.
  - Exposes an HTTP API (`/generate-name`) that accepts JSON payloads with:
    - `application`
    - `role`
    - `environment`
  - Returns a generated hostname in the format:  
    `"{application}-{role}-{environment}-vm{counter}"`

### Custom Terraform Provider

- **Purpose:**  
  Leverages the microservice to generate unique hostnames during Terraform operations.
  
- **Usage Examples:**  
  - **Example 1 (Null Resource):**  
    Uses a null resource with local-exec to call the microservice and an external data source to read the JSON file with the hostname.
  - **Example 2 (Direct Provider):**  
    Uses the custom provider resource that calls the microservice, returning a computed attribute (`generated_name`).

---

## Getting Started

### Prerequisites

- **Vault:**  
  Vault must be installed and configured. The KV V2 secrets engine should be mounted with CAS enabled.
  
- **Terraform:**  
  Ensure Terraform is installed and configured to use the custom provider.
  
- **Docker (optional):**  
  A Dockerfile is provided for containerizing the microservice (see the Dockerfile for details).

### Setup

1. **Vault Configuration:**  
   - Mount the KV secrets engine:
     ```hcl
     resource "vault_mount" "hostnaming_service" {
       path        = "hostnaming"
       type        = "kv"
       options     = { version = "2" }
       description = "providing hostnaming as a service"
     }
     ```
   - Enable CAS on the KV engine (using `vault_kv_secret_backend_v2` resource).

2. **Microservice:**  
   - Configure the environment variables (`VAULT_TOKEN`, `VAULT_ADDR`) and start the Flask microservice.
  
3. **Terraform Examples:**  
   - Review the examples in this repo to see how to integrate both approaches:
     - Null resource + external data source.
     - Direct usage of the custom provider.
  
4. **Run Tests:**  
   - Use the provided scripts (e.g., `run_terraform.sh`) to run performance tests and see metrics.

---

## Disclaimer

This repository is for demonstration purposes only. The code is experimental and is not intended for production use. It is provided "as is" without warranties or guarantees.

---

## Contributions

Feel free to fork this repo and contribute improvements. Please note that changes should be clearly marked as experimental for the PoC.

---

## License

This project is licensed under the MIT License.
