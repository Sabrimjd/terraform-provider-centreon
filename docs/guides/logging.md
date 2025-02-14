---
page_title: "Logging the provider Centreon API"
description: |-
  How to authenticate with Centreon using API V2.
layout: ""
---
# Logging Guide

The Centreon provider supports detailed logging for troubleshooting purposes. Logging can be enabled by setting environment variables before running Terraform commands.

## Environment Variables

- `TF_LOG`: Set the log level. Valid values are (in order of verbosity):
  - `TRACE`
  - `DEBUG`
  - `INFO`
  - `WARN`
  - `ERROR`

- `TF_LOG_PATH`: Specify the log file path. If not set, logs will be written to `terraform-provider-centreon.log` in the current directory.

## Example Usage

```bash
# Enable debug logging
export TF_LOG=DEBUG

# Set custom log file path
export TF_LOG_PATH=/var/log/terraform/centreon-provider.log

# Run Terraform commands
terraform plan
```

## Log File Format

Each log entry includes:
- Timestamp
- Log level
- Message
- Additional context (when available)

Example log entry:
```
2024-01-20T10:15:30Z [INFO] Configuring Centreon client {"server": "centreon.example.com", "protocol": "https", "port": "443", "api_version": "latest"}
```

## Troubleshooting Common Issues

1. **Permission Issues**: Ensure the specified log directory is writable by the user running Terraform.

2. **Missing Logs**: Verify that `TF_LOG` is set to an appropriate level. More detailed logs can be obtained by using `TRACE` level.

3. **Large Log Files**: The provider appends to the log file. Consider rotating or archiving old logs periodically.
