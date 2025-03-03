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
