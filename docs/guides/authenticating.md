---
page_title: "Authenticating to Centreon API"
description: |-
  How to authenticate with Centreon using API V2.
layout: ""
---

Get an API Key with Curl

```bash
curl --request POST \
  --url https://centreon.example.com/centreon/api/v24.10/login \
  --header 'content-type: application/json' \
  --data '{
  "security": {
    "credentials": {
      "login": "username",
      "password": "password"
    }
  }
}'
```
