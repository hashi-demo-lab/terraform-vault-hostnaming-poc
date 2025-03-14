# variable "vms" {
#   description = "Map of VMs with their naming parameters"
#   type = map(object({
#     application = string
#     role        = string
#     environment = string
#   }))
#   default = {
#     vm1 = {
#       application = "splunk"
#       role        = "db"
#       environment = "dev"
#     }
#     vm2 = {
#       application = "splunk"
#       role        = "app"
#       environment = "dev"
#     }
#   }
# }

# # Create a null_resource for each VM that calls the hostnaming service.
# resource "null_resource" "generate_hostname" {
#   for_each = var.vms

#   triggers = {
#     application = each.value.application
#     role        = each.value.role
#     environment = each.value.environment
#   }

#   provisioner "local-exec" {
#     # Write each VM’s hostname output to its own JSON file.
#     command = <<EOT
#       curl -s -X POST "http://localhost:32001/generate-name" \
#       -H "Content-Type: application/json" \
#       -d '{
#           "application": "${each.value.application}",
#           "role": "${each.value.role}",
#           "environment": "${each.value.environment}"
#       }' > ${path.module}/hostname_${each.key}.json
#     EOT
#   }
# }

# # Use an external data source to read each file.
# data "external" "hostname" {
#   for_each = var.vms

#   program = ["bash", "-c", "cat ${path.module}/hostname_${each.key}.json"]

#   depends_on = [null_resource.generate_hostname]
# }

# # Aggregate the hostnames into one output map.
# output "vm_hostnames" {
#   description = "Mapping of VM keys to generated hostnames"
#   value = {
#     for vm, data_source in data.external.hostname :
#     vm => data_source.result["hostname"]
#   }
# }