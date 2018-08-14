package confidant

import (
	"errors"
	"fmt"
)

type Services struct {
	Services []Service `json:"services"`
}

type ServiceResponse struct {
	Result  bool    `json:"result"`
	Service Service `json:"service"`
	Error   string  `json:"error"`
}

type Service struct {
	Enabled          bool          `json:"enabled"`
	ID               string        `json:"id"`
	Revision         int           `json:"revision"`
	Credentials      []*Credential `json:"credentials"`
	BlindCredentials []*Credential `json:"blind_credentials"`
	Account          string        `json:"account"`
	Error            string        `json:"error"`
	ModifiedBy       string        `json:"modified_by"`
	ModifiedDate     string        `json:"modified_date"`
}

// GetServices fetches the list of services
// It returns a pointer to a Services struct
func (c *Client) GetServices() (*Services, error) {
	var services Services
	err := c.Request("GET", "/v1/services", nil, &services)
	if err != nil {
		return nil, fmt.Errorf("Got an error when making the Confidant request: %e", err)
	}
	return &services, nil
}

// GetService fetches details for a service.
// It returns a pointer to a Service struct.
func (c *Client) GetService(serviceName string) (*Service, error) {
	if service, ok := c.services[serviceName]; ok {
		return service, nil
	}
	var service Service
	err := c.Request("GET", "/v1/services/"+serviceName, nil, &service)
	if err != nil {
		if err.Error() == "NotFound" {
			return nil, errors.New("Service Doesn't Exist")
		} else {
			return nil, err
		}
	} else if service.Error != "" {
		return nil, errors.New(service.Error)
	}
	c.services[serviceName] = &service
	return &service, nil
}

// CreateService creates a new service.
// It returns a pointer to a Service struct.
func (c *Client) CreateService(serviceName string, credentialNames []string) (*Service, error) {
	service, err := c.GetService(serviceName)
	if err == nil {
		err = c.EnsureGrants(serviceName)
		if err != nil {
			return service, fmt.Errorf("Service already exists, but got an error when trying to ensure grants: %e", err)
		}
		return service, errors.New("Service Already Exists")
	} else if err.Error() != "Service Doesn't Exist" {
		return nil, err
	}
	err = c.CheckRole(serviceName)
	if err != nil {
		return nil, err
	}
	credentials, err := c.FindCredentialsByName(credentialNames)
	if err != nil {
		return nil, err
	}
	credentialIDs := make([]string, len(credentials))
	for i, v := range credentials {
		credentialIDs[i] = v.ID
	}
	body := RequestBody{
		Credentials: credentialIDs,
		Enabled:     true,
	}
	var response ServiceResponse
	err = c.Request("PUT", "/v1/services/"+serviceName, &body, &response)
	if err != nil {
		return nil, err
	} else if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	if response.Service.Revision != 0 {
		err := c.EnsureGrants(serviceName)
		if err != nil {
			return nil, err
		}
	}
	c.services[serviceName] = &response.Service
	return &response.Service, nil
}

// SetServiceCredentials sets the credentials for an existing service.
// It returns a pointer to a Service struct.
func (c *Client) SetServiceCredentials(serviceName string, credentialNames []string) (*Service, error) {
	service, err := c.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	err = c.EnsureGrants(serviceName)
	if err != nil {
		return nil, err
	}
	credentials, err := c.FindCredentialsByName(credentialNames)
	if err != nil {
		return nil, err
	}
	credentialIDs := make([]string, 0, len(credentials))
	for _, credential := range credentials {
		credentialIDs = append(credentialIDs, credential.ID)
	}

	body := RequestBody{
		Credentials:      credentialIDs,
		BlindCredentials: []string{},
		Account:          service.Account,
		Enabled:          service.Enabled,
	}
	var response Service
	err = c.Request("PUT", "/v1/services/"+serviceName, &body, &response)
	if err != nil {
		return nil, err
	} else if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	if response.Revision != 0 {
		err := c.EnsureGrants(serviceName)
		if err != nil {
			return nil, fmt.Errorf("Could not ensure grants: %e", err)
		}
	}
	c.services[serviceName] = &response
	return &response, nil
}

// UpdateServiceCredentials updates an existing service by adding or removing credentials.
// It returns a pointer to a Service struct.
func (c *Client) UpdateServiceCredentials(serviceName string, addCredentialNames []string, removeCredentialNames []string) (*Service, error) {
	service, err := c.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	err = c.EnsureGrants(serviceName)
	if err != nil {
		return nil, err
	}
	addCredentials, err := c.FindCredentialsByName(addCredentialNames)
	if err != nil {
		return nil, err
	}
	removeCredentials, err := c.FindCredentialsByName(removeCredentialNames)
	if err != nil {
		return nil, err
	}
	merged := createCredentialMap(service.Credentials, addCredentials, removeCredentials)
	credentialIDs := getCredentialIDs(merged)

	body := RequestBody{
		Credentials:      credentialIDs,
		BlindCredentials: []string{},
		Account:          service.Account,
		Enabled:          service.Enabled,
	}
	var response Service
	err = c.Request("PUT", "/v1/services/"+serviceName, &body, &response)
	if err != nil {
		return nil, err
	} else if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	if response.Revision != 0 {
		err := c.EnsureGrants(serviceName)
		if err != nil {
			return nil, fmt.Errorf("Could not ensure grants: %e", err)
		}
	}
	c.services[serviceName] = &response
	return &response, nil
}

// EnableService updates an existing service by setting enabled to true
// It returns a pointer to a Service struct.
func (c *Client) EnableService(serviceName string) (*Service, error) {
	service, err := c.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	err = c.EnsureGrants(serviceName)
	if err != nil {
		return nil, err
	}
	credentialIDs := getCredentialIDs(service.Credentials)

	body := RequestBody{
		Credentials:      credentialIDs,
		BlindCredentials: []string{},
		Account:          service.Account,
		Enabled:          true,
	}
	var response Service
	err = c.Request("PUT", "/v1/services/"+serviceName, &body, &response)
	if err != nil {
		return nil, err
	} else if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	if response.Revision != 0 {
		err := c.EnsureGrants(serviceName)
		if err != nil {
			return nil, err
		}
	}
	c.services[serviceName] = &response
	return &response, nil
}

// DisableService updates an existing service by setting Enabled to false
// and Credentials to empty.
// It returns a pointer to a Service struct.
func (c *Client) DisableService(serviceName string) (*Service, error) {
	service, err := c.GetService(serviceName)
	if err != nil {
		return nil, err
	}

	body := RequestBody{
		Credentials:      []string{},
		BlindCredentials: []string{},
		Account:          service.Account,
		Enabled:          false,
	}
	var response Service
	err = c.Request("PUT", "/v1/services/"+serviceName, &body, &response)
	if err != nil {
		return nil, err
	} else if response.Error != "" {
		return nil, errors.New(response.Error)
	}
	c.services[serviceName] = &response
	return &response, nil
}
