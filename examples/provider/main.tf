terraform {
  required_providers {
    centreon = {
      source  = "registry.terraform.io/smjed/centreon"
      version = "0.1.0"
    }
  }
}

provider "centreon" {
  protocol    = "https"
  server      = "centreon.acme.lan"
  port        = "443"
  api_version = "latest"
  api_key     = "YOUR_API_KEY"
}

data "centreon_platform_info" "info" {
}

output "installation_status" {
  value = data.centreon_platform_info.info.is_installed
}

output "upgrade_available" {
  value = data.centreon_platform_info.info.has_upgrade_available
}
