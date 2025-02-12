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

# Get all host groups
data "centreon_host_groups" "all" {
  limit = 50
  page  = 1
}

# Search for specific host groups
data "centreon_host_groups" "linux_prod" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "LINUX"
  }
}

data "centreon_host_groups" "prod" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "PROD"
  }
}

# Output all groups
output "all_groups" {
  value = [for g in data.centreon_host_groups.all.groups : {
    id   = g.id
    name = g.name
  }]
}

# Output Linux production groups
output "linux_groups" {
  value = data.centreon_host_groups.linux_prod.groups
}

# Example: Use group information in a host resource
resource "centreon_host" "linux_server" {
  monitoring_server_id = 1
  name                 = "linux-server-01"
  address              = "192.168.1.100"
  alias                = "Linux Production Server"

  # Assign to both LINUX and PROD groups
  groups = concat(
    [for g in data.centreon_host_groups.linux_prod.groups : g.id],
    [for g in data.centreon_host_groups.prod.groups : g.id]
  )

  is_activated = true
}
