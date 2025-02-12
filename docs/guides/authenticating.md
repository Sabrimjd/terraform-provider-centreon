---
page_title: "Authenticating to Centreon API"
---

Get an API Key with Curl

```bash
curl --request POST \
  --url https://centreon.example.com/centreon/api/v24.10/login \
  --header 'content-type: application/json' \
  --data '{
  "security": {
    "credentials": {
      "login": "{{username}}",
      "password": "{{password}}"
    }
  }
}'
```
