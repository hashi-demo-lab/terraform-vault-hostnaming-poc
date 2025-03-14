terraform {
  required_providers {
    naming = {
      source = "local/terraform-providers/naming"
      version = "1.0.0"
    }
  }
}

resource "naming" "test" {
  workspace_id = "ws-12345"
  vms = {
    "vm1" = {}
  }
}

output "generated_names" {
  value = naming.test.generated_names
} 