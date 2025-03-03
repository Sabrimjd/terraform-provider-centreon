# Search for Linux Oracle templates
data "centreon_host_templates" "oracle_templates" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "Linux_Oracle_9"
  }
}
