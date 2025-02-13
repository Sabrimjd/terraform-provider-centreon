terraform {
  required_providers {
    centreon = {
      source = "github.com/Sabrimjd/centreon"
    }
  }
  required_version = ">= 1.0.0"
}

provider "centreon" {
  protocol    = "https"
  server      = "centreon.example.com"
  port        = "443"
  api_version = "latest"
  api_key     = "your-api-key"
}

# Example of creating a basic web server host
resource "centreon_host" "web_server" {
  monitoring_server_id = 1
  name                 = "web-server-01"
  address              = "192.168.1.100"
  alias                = "Production Web Server"

  # Basic monitoring configuration
  check_command_id      = 1             # Assuming 1 is the ID for HTTP check
  check_command_args    = ["80", "300"] # Port 80, timeout 300s
  max_check_attempts    = 3
  normal_check_interval = 5 # Check every 5 minutes when OK
  retry_check_interval  = 1 # Check every minute when not OK

  # Notification configuration
  notification_enabled       = 1  # Enable notifications
  notification_options       = 5  # DOWN (1) + RECOVERY (4)
  notification_interval      = 30 # Notify every 30 minutes
  notification_timeperiod_id = 1  # Assuming 1 is 24x7 timeperiod

  # Group assignments
  templates  = [1] # Basic host template
  groups     = [1] # Web servers group
  categories = [1] # Production category

  # Custom macros
  macros = [
    {
      name        = "HTTP_PORT"
      value       = "80"
      is_password = false
      description = "Web server port"
    },
    {
      name        = "ENVIRONMENT"
      value       = "production"
      is_password = false
      description = "Environment type"
    }
  ]

  # Additional metadata
  comment      = "Main production web server"
  is_activated = true
}

# Example of creating a database server with SNMP monitoring
resource "centreon_host" "db_server" {
  monitoring_server_id = 1
  name                 = "db-server-01"
  address              = "192.168.1.101"
  alias                = "Production Database Server"

  # SNMP configuration
  snmp_community = "public"
  snmp_version   = "2c"

  # Basic monitoring configuration
  check_command_id      = 2 # Assuming 2 is the ID for MySQL check
  max_check_attempts    = 4
  normal_check_interval = 3
  retry_check_interval  = 1

  # Group assignments
  templates  = [2] # Database server template
  groups     = [2] # Database servers group
  categories = [1] # Production category

  # Custom macros for database monitoring
  macros = [
    {
      name        = "MYSQL_PORT"
      value       = "3306"
      is_password = false
      description = "MySQL server port"
    },
    {
      name        = "MYSQL_USER"
      value       = "monitoring"
      is_password = false
      description = "MySQL monitoring user"
    },
    {
      name        = "MYSQL_PASSWORD"
      value       = "secret123"
      is_password = true
      description = "MySQL monitoring user password"
    }
  ]

  is_activated = true
}
