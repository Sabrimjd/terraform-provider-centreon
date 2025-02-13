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

# Get all host templates
data "centreon_host_templates" "all" {
  limit = 200
  page  = 1
}

# Search for specific templates
data "centreon_host_templates" "elk_templates" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "Linux_Centos_7_extension_ELK_PROD"
  }
}

# Search for Linux Oracle templates
data "centreon_host_templates" "oracle_templates" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "Linux_Oracle_9"
  }
}

# Output all templates basic info
output "all_templates" {
  value = [for t in data.centreon_host_templates.all.templates : {
    id    = t.id
    name  = t.name
    alias = t.alias
  }]
}

# Output ELK template details
output "elk_template_details" {
  value = data.centreon_host_templates.elk_templates.templates
}

# Example: Use template information in a host resource
resource "centreon_host" "elastic_server" {
  monitoring_server_id = 1
  name                 = "elastic-server-01"
  address              = "192.168.1.100"
  alias                = "Elastic Search Server"

  # Combine ELK and Oracle templates
  templates = concat(
    [for t in data.centreon_host_templates.elk_templates.templates : t.id],
    [for t in data.centreon_host_templates.oracle_templates.templates : t.id]
  )

  # Set monitoring parameters based on template defaults
  active_check_enabled   = 2 # Use template default
  notification_enabled   = 2 # Use template default
  flap_detection_enabled = 2 # Use template default

  is_activated = true
}
