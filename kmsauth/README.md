# KMS Authentication
Confidant supports authentication via KMS. This library generates and can validate authentication tokens, based on https://github.com/lyft/python-kmsauth

The tokens are generated in v2 format, which looks like:

* username: "2/user/terraform-provider-confidant"
* encryption context: {"to":"confidant-production","from":"terraform-provider-confidant","user_type":"user"}

## Usage
### Generating username and token
Decrypting tokens requires the username and the token, so when passing this to a service, you should pass both along.

```go
package main

import "github.com/stripe/go-confidant-client/kmsauth"

func main() {
  // KMS key to use for authentication
  key := "alias/authnz-production"
  // The service being authenticated
  to := "confidant-production"
  // The user for whom the token is being generated
  from := "terraform-provider-confidant"
  userType := "user"
  region := "us-east-1"
  generator := kmsauth.NewTokenGenerator(key, to, from, userType, region)
  username := generator.GetUsername()
  token := generator.GetToken()
}
```
