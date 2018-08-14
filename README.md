# go-confidant [![Build Status](https://travis-ci.com/stripe/go-confidant-client.svg?branch=master)](https://travis-ci.com/stripe/go-confidant-client)
A Go client library for [Confidant](https://lyft.github.io/confidant/).

## Installation
`$ go get github.com/stripe/go-confidant-client`

## Usage
### Initializing the client
Creating a client requires a url, a http client and a KMS auth token generator.

`kmsauth.NewTokenGenerator()` takes the id of the Confidant KMS key ("authkey"), the Confidant IAM role ("to"), "from" (the user making the change, this could be an IAM role or AWS username[1]), user type (either "user" or "service") and the region in which the AWS KMS encrypt call will be made.

1: It's possible to restrict access to Confidant by adding an IAM policy to the Confidant KMS key specifying that the `from` field should match the username or IAM role of the person making the AWS request. You can read more about how Confidant uses KMS for authentication [here](https://medium.com/@arpith/how-confidant-uses-kms-for-authentication-4aa14d5f6b91).

```go
import (
	"github.com/stripe/go-confidant-client/kmsauth"
	"net/http"
)

func initClient() *Client {
	authkey := "key"
	to := "ConfidantServer"
	from := "username"
	userType := "user"
	region := "us-east-1"
	url := "confidant-url"
	httpClient := &http.Client{}
	generator := kmsauth.NewTokenGenerator(authkey, to, from, userType, region)
	c := NewClient(url, httpClient, &generator)
	return &c
}
```
### Services
#### Get Services
To get a list of services call `client.GetServices()`.
```go
func ExampleGetServices() {
	c := initClient()
	services, err := c.GetServices()
	if err != nil {
		log.Printf("Got an error when getting services: %e", err)
	}
	fmt.Printf("%+v\n", services)
}
```

#### Get Service
To fetch a service, pass the service name to `client.GetService()`.
```go
func ExampleGetService() {
	name := "name"
	c := initClient()
	service, err := c.GetService(name)
	if err != nil {
		log.Printf("Got an error when getting the service named %s: %e", name, err)
	}
	fmt.Printf("%+v\n", service)
}
```

The Service type that is returned looks like:

```go
type Service struct {
	Enabled          bool          `json:"enabled"`
	ID               string        `json:"id"`
	Revision         string        `json:"string"`
	Credentials      []*Credential `json:"credentials"`
	BlindCredentials []*Credential `json:"blind_credentials"`
	Account          string        `json:"account"`
	Error            string        `json:"error"`
}
```

#### Create a Service
To create a service, pass the service name and a slice of credential names to `client.CreateService()`.
```go
func ExampleCreateService() {
	name := "test-go-confidant-create-service"
	c := initClient()
	service, err := c.CreateService(name, []string{"test-credential"})
	if err != nil {
		log.Printf("Got an error when creating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}
```

#### Update a Service
##### Set a Service's Credentials
To set a service's credentials, pass the service name and a slice of credential names to `client.SetServiceCredentials()`. This will overwrite the list of credentials currently assigned to the service. To update non-destructively, use [`client.UpdateServiceCredentials()`](https://github.com/stripe/go-confidant-client#update-a-services-credentials) instead.

```go
func ExampleSetServiceCredentials() {
	name := "service-name"
	c := initClient()
	service, err := c.SetServiceCredentials(name, []string{"test-credential"})
	if err != nil {
		log.Printf("Got an error when setting credentials for the service named %s: %e", name, err)
	}
	fmt.Println(service)
}
```

##### Update a Service's Credentials
To update a service's credentials, pass the service name, a slice of credential names to add and a slice of credential names to remove to `client.UpdateServiceCredentials()`.

```go
func ExampleUpdateServiceCredentials() {
	name := "service-name"
	c := initClient()
	service, err := c.UpdateServiceCredentials(name, []string{"test-credential"}, []string{})
	if err != nil {
		log.Printf("Got an error when updating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}
```

##### Enable a Service
To enable a service, pass the service name to `client.EnableService()`.
```go
func ExampleEnableService() {
	name := "service-name"
	c := initClient()
	service, err := c.EnableService(name)
	if err != nil {
		log.Printf("Got an error when updating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}
```
##### Disable a Service
To disable a service, pass the service name to `client.DisableService()`.
```go
func ExampleDisableService() {
	name := "service-name"
	c := initClient()
	service, err := c.DisableService(name)
	if err != nil {
		log.Printf("Got an error when updating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}
```

### Credentials
#### Assign Credentials
To assign a credential to a service, pass the service name and credential name to `client.AssignCredential()`. This updates the service and adds the credential.
```go
func ExampleAssignCredential() {
	credential := "yet-another-test-credential"
	service := "name"
	c := initClient()
	err := c.AssignCredential(service, credential)
	if err != nil {
		log.Printf("Got an error when assigning %s to %s", credential, service)
	}
	fmt.Println(c.services)
}
```

#### Unassign Credentials
To remove a credential from a service, pass the service name and credential name to `client.UnassignCredential()`. This updates the service and removes the credential.
```go
func ExampleUnassignCredential() {
	credential := "yet-another-test-credential"
	service := "name"
	c := initClient()
	err := c.UnassignCredential(service, credential)
	if err != nil {
		log.Printf("Got an error when unassigning %s to %s", credential, service)
	}
	fmt.Println(c.services)
}
```

#### Find Credentials By Name
To get credentials by name, pass the names to `client.FindCredentialsByName()`. An error is returned if any of the names passed in do not match credentials.
```go
func ExampleFindCredentialsByName() {
	name := "yet-another-test-credential"
	c := initClient()
	credentials, err := c.FindCredentialsByName([]string{name})
	if err != nil {
		log.Printf("Got an error when finding credentials named %s", name)
	}
	fmt.Printf("%+v", credentials[0])
}
```


### Grants
To make sure a service has grants to encrypt and decrypt, pass the service name to `client.EnsureGrants()`. The grants can be checked by calling `client.GetGrants()` with the service name.
```go
func ExampleEnsureGrants() {
	name := "service-name"
	c := initClient()
	err := c.EnsureGrants(name)
	if err != nil {
		log.Printf("Got an error when ensuring grants for service named %s: %e", name, err)
	}
	grants, err := c.GetGrants(name)
	if err != nil {
		log.Printf("Got an error when getting grants for service named %s: %e", name, err)
	}
	fmt.Println(grants)
}
```

### Role
To check that a service name is a valid IAM role, call `client.CheckRole()` with the service name.
```go
func ExampleCheckRole() {
	name := "service-name"
	c := initClient()
	err := c.CheckRole(name)
	if err != nil {
		log.Printf("Got an error when checking if %s is a valid IAM role: %e", name, err)
	}
}
```


