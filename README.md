# Terraform Provider for Centreon

This Terraform Provider allows you to interact with Centreon through its API. It provides the ability to manage and query Centreon resources through Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
```sh
git clone git@github.com:your-username/terraform-provider-centreon.git
```

2. Enter the repository directory
```sh
cd terraform-provider-centreon
```

3. Build the provider
```sh
go build -o terraform-provider-centreon
```

## Documentation

The provider documentation can be found in the [docs](docs/) directory.

## Examples

Examples of how to use the provider can be found in the [examples](examples/) directory.

## Using the provider

```hcl
terraform {
  required_providers {
    centreon = {
      source = "Sabrimjd/centreon"
    }
  }
}

provider "centreon" {
  protocol     = "https"
  server       = "centreon.example.com"
  port         = "443"
  api_version  = "latest"
  api_key      = "your-api-key"
}

# Get platform information
data "centreon_platform_info" "info" {
}

output "is_installed" {
  value = data.centreon_platform_info.info.is_installed
}

output "has_upgrade_available" {
  value = data.centreon_platform_info.info.has_upgrade_available
}

# Search for specific hosts
data "centreon_host_search" "elastic_hosts" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "HPLESMBE1-010"
  }
}

output "host_details" {
  value = data.centreon_host_search.elastic_hosts.hosts
}

# Search for hosts in configuration
data "centreon_configuration_hosts" "elastic_hosts" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "HPLESMBE1-010"
  }
}

output "hosts" {
  value = data.centreon_configuration_hosts.elastic_hosts.hosts
}

# Access specific host attributes
output "first_host_name" {
  value = length(data.centreon_configuration_hosts.elastic_hosts.hosts) > 0 ? data.centreon_configuration_hosts.elastic_hosts.hosts[0].name : ""
}

output "first_host_groups" {
  value = length(data.centreon_host_search.elastic_hosts.hosts) > 0 ? data.centreon_host_search.elastic_hosts.hosts[0].groups : []
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build`. This will build the provider and put the provider binary in the current directory.

```sh
go build -o terraform-provider-centreon
```

To generate or update documentation, use the following commands:

```sh
go generate ./...
```

In order to run the full suite of tests, run `make test`.

```sh
make test
```

## License

[MPL-2.0](LICENSE)