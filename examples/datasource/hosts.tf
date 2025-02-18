# Search for specific hosts
data "centreon_hosts" "some_hosts" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "HPLESMBE1-010"
  }
}
