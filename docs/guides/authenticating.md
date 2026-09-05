---
page_title: "Authenticating to Centreon API"
description: |-
  How to authenticate with Centreon using API V2.
layout: ""
---

# Authenticating to Centreon API

## Getting an API Key

You can obtain an API key from Centreon using the following curl command:

```bash
curl --request POST \\
  --url https://centreon.example.com/centreon/api/latest/login \\
  --header 'content-type: application/json' \\
  --data '{
  "login": "username",
  "password": "password"
}
```

you can now use this api key in the provider under the api_key parameters.

## Using Vault and Variables

To avoid storing plain text API keys in your Terraform files, you can use Vault and a Terraform variable. Here is an example of how to retrieve the API key from Vault and pass it to Terraform:

```bash
export TF_VAR_centreon_api_key=$(vault kv get -field=api_key secret/centreon)
```

Then, you can reference the variable in your Terraform provider configuration:

```hcl
provider "centreon" {
  api_key = var.centreon_api_key
}
```

Make sure to define the `centreon_api_key` variable in your Terraform configuration:

```hcl
variable "centreon_api_key" {
  description = "API key for Centreon"
  type        = string
  sensitive   = true
}
```
