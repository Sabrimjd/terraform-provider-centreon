# Centreon Provider Examples

This directory contains examples of using the Centreon Terraform provider.

## Example Files

- `provider-install-verification/`: Basic example that verifies the provider installation and shows how to fetch platform information.

## Running the Examples

1. First, ensure you have valid Centreon API credentials.

2. Update the provider configuration in the example with your Centreon server details:
```hcl
provider "centreon" {
  protocol     = "https"
  server       = "your-centreon-server"
  port         = "443"
  api_version  = "latest"
  api_key      = "your-api-key"
}
```

3. Initialize Terraform:
```sh
terraform init
```

4. Run Terraform:
```sh
terraform apply
```

This will execute the example configuration and show you the outputs.

## Notes

- Make sure to replace sensitive values (like api_key) with your actual credentials
- Never commit sensitive credentials to version control
- Consider using variables or environment variables for sensitive values in production
