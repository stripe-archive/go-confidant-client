package confidant

import (
	"fmt"
	"log"
)

func ExampleRequest() {
	var services Services
	c := initClient()
	err := c.Request("GET", "/v1/services", nil, &services)
	if err != nil {
		log.Printf("Got an error when making the Confidant request: %e", err)
	}
	fmt.Println(services)
}

func ExampleGetServices() {
	c := initClient()
	services, err := c.GetServices()
	if err != nil {
		log.Printf("Got an error when getting services: %e", err)
	}
	fmt.Printf("%+v\n", services)
}

func ExampleGetService() {
	name := "name"
	c := initClient()
	service, err := c.GetService(name)
	if err != nil {
		log.Printf("Got an error when getting the service named %s: %e", name, err)
	}
	fmt.Printf("%+v\n", service)
}

func ExampleCreateService() {
	name := "test-go-confidant-create-service"
	c := initClient()
	service, err := c.CreateService(name, []string{"test-credential"})
	if err != nil {
		log.Printf("Got an error when creating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}

func ExampleSetServiceCredentials() {
	name := "service-name"
	c := initClient()
	service, err := c.SetServiceCredentials(name, []string{"test-credential"})
	if err != nil {
		log.Printf("Got an error when setting credentials for the service named %s: %e", name, err)
	}
	fmt.Println(service)
}

func ExampleUpdateServiceCredentials() {
	name := "service-name"
	c := initClient()
	service, err := c.UpdateServiceCredentials(name, []string{"test-credential"}, []string{})
	if err != nil {
		log.Printf("Got an error when updating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}

func ExampleEnableService() {
	name := "service-name"
	c := initClient()
	service, err := c.EnableService(name)
	if err != nil {
		log.Printf("Got an error when updating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}

func ExampleDisableService() {
	name := "service-name"
	c := initClient()
	service, err := c.DisableService(name)
	if err != nil {
		log.Printf("Got an error when updating the service named %s: %e", name, err)
	}
	fmt.Println(service)
}

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

func ExampleCheckRole() {
	name := "service-name"
	c := initClient()
	err := c.CheckRole(name)
	if err != nil {
		log.Printf("Got an error when checking if %s is a valid IAM role: %e", name, err)
	}
}

func ExampleFindCredentialsByName() {
	name := "yet-another-test-credential"
	c := initClient()
	credentials, err := c.FindCredentialsByName([]string{name})
	if err != nil {
		log.Printf("Got an error when finding credentials named %s", name)
	}
	fmt.Printf("%+v", credentials[0])
}

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
