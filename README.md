# Terraform Provider for Centreon

This Terraform Provider allows you to interact with Centreon through its API. It provides the ability to manage and query Centreon resources through Terraform.

> ⚠️ **Warning**: This provider is in early stages of development and is not ready for production use. Features may be incomplete, and breaking changes can occur without notice. Use it for testing and evaluation purposes only.

You can see the official documentation on the Terraform Provider dedicated webpage here : [https://registry.terraform.io/providers/Sabrimjd/centreon/latest/docs](https://registry.terraform.io/providers/Sabrimjd/centreon/latest/docs)

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
  protocol                         = "https"
  server                           = "centreon.example.com"
  port                             = "443"
  api_version                      = "latest"
  api_key                          = "your-api-key"
  generate_and_reload_configuration = true
}

# Create a new host
resource "centreon_host" "web_server" {
  monitoring_server_id = 1
  name                = "web-server-01"
  address             = "192.168.1.100"
  alias               = "Production Web Server"

  # Server checks configuration
  check_command_id       = 1
  check_command_args     = ["80", "300"]
  max_check_attempts     = 3
  normal_check_interval  = 5
  retry_check_interval   = 1

  # Notification settings
  notification_enabled     = 1
  notification_options     = 5  # DOWN (1) + RECOVERY (4)
  notification_interval    = 30
  notification_timeperiod_id = 1

  # Host grouping
  templates  = [1]  # List of template IDs
  groups     = [1]  # List of group IDs

  # Activation state
  is_activated = true
}

# Get platform information
data "centreon_platform_info" "info" {
}

output "is_installed" {
  value = data.centreon_platform_info.info.is_installed
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
