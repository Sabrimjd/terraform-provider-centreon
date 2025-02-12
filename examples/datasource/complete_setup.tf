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

# Get central monitoring server
data "centreon_monitoring_servers" "central" {
  limit = 1
  page  = 1
  search = {
    name  = "name"
    value = "Central"
  }
}

# Get all required templates
data "centreon_host_templates" "elk_stack" {
  limit = 100
  page  = 1
  search = {
    name  = "name"
    value = "Linux_Centos_7_extension_ELK_PROD"
  }
}

data "centreon_host_templates" "oracle" {
  limit = 100
  page  = 1
  search = {
    name  = "name"
    value = "Linux_Oracle_9"
  }
}

# Get all required groups
data "centreon_host_groups" "linux" {
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

data "centreon_host_groups" "elastic" {
  limit = 10
  page  = 1
  search = {
    name  = "name"
    value = "Elastic"
  }
}

# Create multiple Elastic Stack servers with consistent configuration
locals {
  elastic_servers = {
    "es01" = {
      name    = "elastic-01"
      address = "10.0.1.10"
      alias   = "Elastic Search Node 1"
    },
    "es02" = {
      name    = "elastic-02"
      address = "10.0.1.11"
      alias   = "Elastic Search Node 2"
    },
    "kibana" = {
      name    = "kibana-01"
      address = "10.0.1.20"
      alias   = "Kibana Dashboard"
    }
  }
}

# Create the Elastic Stack servers
resource "centreon_host" "elastic_cluster" {
  for_each = local.elastic_servers

  monitoring_server_id = data.centreon_monitoring_servers.central.servers[0].id
  name                 = each.value.name
  address              = each.value.address
  alias                = each.value.alias

  # Apply all relevant templates
  templates = concat(
    [for t in data.centreon_host_templates.elk_stack.templates : t.id],
    [for t in data.centreon_host_templates.oracle.templates : t.id]
  )

  # Apply all relevant groups
  groups = concat(
    [for g in data.centreon_host_groups.linux.groups : g.id],
    [for g in data.centreon_host_groups.prod.groups : g.id],
    [for g in data.centreon_host_groups.elastic.groups : g.id]
  )

  # Common monitoring settings
  active_check_enabled   = 2 # Use template default
  passive_check_enabled  = 2 # Use template default
  notification_enabled   = 1 # Enable notifications
  notification_options   = 5 # DOWN (1) + RECOVERY (4)
  notification_interval  = 30
  max_check_attempts     = 3
  normal_check_interval  = 5
  retry_check_interval   = 1
  freshness_checked      = 2 # Use template default
  flap_detection_enabled = 2 # Use template default
  event_handler_enabled  = 2 # Use template default

  # Custom macros for Elastic Stack monitoring
  macros = [
    {
      name        = "_ES_PORT"
      value       = "9200"
      is_password = false
      description = "Elasticsearch HTTP port"
    },
    {
      name        = "_CLUSTER_NAME"
      value       = "production"
      is_password = false
      description = "Elasticsearch cluster name"
    }
  ]

  is_activated = true
}

# Verify the setup with outputs
output "monitoring_server" {
  value = data.centreon_monitoring_servers.central.servers[0]
}

output "applied_templates" {
  value = distinct(concat(
    [for t in data.centreon_host_templates.elk_stack.templates : t.name],
    [for t in data.centreon_host_templates.oracle.templates : t.name]
  ))
}

output "applied_groups" {
  value = distinct(concat(
    [for g in data.centreon_host_groups.linux.groups : g.name],
    [for g in data.centreon_host_groups.prod.groups : g.name],
    [for g in data.centreon_host_groups.elastic.groups : g.name]
  ))
}

output "elastic_hosts" {
  value = {
    for k, host in centreon_host.elastic_cluster : k => {
      name    = host.name
      address = host.address
      alias   = host.alias
    }
  }
}
