# Get platform information
data "centreon_platform_info" "info" {
}

output "is_installed" {
  value = data.centreon_platform_info.info.is_installed
}