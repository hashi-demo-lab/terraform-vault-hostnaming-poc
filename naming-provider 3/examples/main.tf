terraform {
  required_providers {
    naming = {
      source = "local/terraform-providers/naming"
    }
  }
}

resource "naming" "example" {
  workspace_id = "ws-12345"
  vms = {
    "vm1" = ""
    "vm2" = ""
    "vm3" = ""
  }
}

output "generated_names" {
  value = naming.example.generated_names
} 