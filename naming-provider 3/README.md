# Terraform Provider for VM Naming

This Terraform provider allows you to generate names for VMs using a naming service API.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```sh
go install
```

## Using the provider

```hcl
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
  }
}

output "generated_names" {
  value = naming.example.generated_names
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ go install
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
``` 