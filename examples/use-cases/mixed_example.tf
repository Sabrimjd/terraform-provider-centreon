terraform {
  required_providers {
    centreon = {
      source = "smjed.net/terraform-providers/centreon"
    }
  }
}

provider "centreon" {
  protocol    = "https"
  server      = "your-centreon-server"
  port        = "443"
  api_version = "latest"
  api_key     = "your-api-key"
}

# Get monitoring server information
data "centreon_monitoring_servers" "servers" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "Central" # Looking for the central server
  }
}

# Get host templates
data "centreon_host_templates" "elk_templates" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "Linux_Centos_7_extension_ELK_PROD"
  }
}

# Get host groups
data "centreon_host_groups" "prod_groups" {
  limit = 50
  page  = 1
  search = {
    name  = "name"
    value = "PROD"
  }
}

# Create a host using the data from our data sources
resource "centreon_host" "elastic_server" {
  monitoring_server_id = data.centreon_monitoring_servers.servers.servers[0].id
  name                 = "HPLESMBE1-020"
  address              = "10.46.120.20"
  alias                = "Elastic Stack UAT"

  # Template assignments using data source
  templates = [
    data.centreon_host_templates.elk_templates.templates[0].id
  ]

  # Group assignments using data source
  groups = [
    data.centreon_host_groups.prod_groups.groups[0].id
  ]

  is_activated = true
}

# Output the details for verification
output "monitoring_server" {
  value = data.centreon_monitoring_servers.servers.servers[0]
}

output "template_details" {
  value = data.centreon_host_templates.elk_templates.templates[0]
}

output "group_details" {
  value = data.centreon_host_groups.prod_groups.groups
}
