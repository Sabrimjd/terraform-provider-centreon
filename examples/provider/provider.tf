provider "centreon" {
  protocol                          = "https"
  server                            = "centreon.acme.lan"
  port                              = "443"
  api_version                       = "latest"
  api_key                           = "YOUR_API_KEY"
  generate_and_reload_configuration = true
}
