terraform {
  required_providers {
    naming = {
      source = "local/terraform-providers/naming"
    }
  }
}

variable "vms" {
  description = "Map of VMs with their naming parameters"
  type = map(object({
    application = string
    role        = string
    environment = string
  }))
  default = {
    vm1 = {
      application = "prometheus"
      role        = "app"
      environment = "prod"
    }
    vm2 = {
      application = "prometheus"
      role        = "web"
      environment = "prod"
    }
    vm3 = {
      application = "prometheus"
      role        = "web"
      environment = "prod"
    }
  }
}

resource "naming" "example" {
  for_each    = var.vms

  application = each.value.application
  role        = each.value.role
  environment = each.value.environment
}

output "generated_hostnames" {
  value = { for key, res in naming.example : key => res.generated_name }
}
