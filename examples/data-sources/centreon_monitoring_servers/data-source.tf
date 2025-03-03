# Search for a specific monitoring server
data "centreon_monitoring_servers" "central" {
  limit = 1
  page  = 1
  search = {
    name  = "name"
    value = "Central"
  }
}
