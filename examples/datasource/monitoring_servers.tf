terraform {
  required_providers {
    centreon = {
      source = "smjed.net/terraform-providers/centreon"
    }
  }
}

provider "centreon" {
  protocol    = "https"
  server      = "centreon.example.com"
  port        = "443"
  api_version = "latest"
  api_key     = "your-api-key"
}

# Get all monitoring servers
data "centreon_monitoring_servers" "all" {
  limit = 100
  page  = 1
}

# Search for a specific monitoring server
data "centreon_monitoring_servers" "central" {
  limit = 1
  page  = 1
  search = {
    name  = "name"
    value = "Central"
  }
}

# Output basic server information
output "all_servers" {
  value = [for s in data.centreon_monitoring_servers.all.servers : {
    id      = s.id
    name    = s.name
    address = s.address
  }]
}

# Output detailed information about the Central server
output "central_server" {
  value = data.centreon_monitoring_servers.central.servers[0]
}

# Example: Use server information in a host resource
resource "centreon_host" "example" {
  monitoring_server_id = data.centreon_monitoring_servers.central.servers[0].id
  name                 = "example-host"
  address              = "192.168.1.100"
  alias                = "Example Host"
  is_activated         = true
}
