package confidant

import (
	"fmt"
	"log"
	"time"
)

type Grants struct {
	EncryptGrant bool `json:"encrypt_grant"`
	DecryptGrant bool `json:"decrypt_grant"`
}

type GrantsResponse struct {
	Grants Grants `json:"grants"`
	Error  string `json:"error"`
}

// GetGrants fetches a service's grants.
// It makes a GET request to /v1/grants/serviceName.
func (c *Client) GetGrants(serviceName string) (*Grants, error) {
	var response GrantsResponse
	err := c.Request("GET", "/v1/grants/"+serviceName, nil, &response)
	return &response.Grants, err
}

// EnsureGrants adds encrypt and decrypt grants for a service.
// It makes 10 PUT requests to /v1/grants/serviceName.
// If the error from Confidant indicates that repeated requests will not succeed,
// it returns an error immediately.
// If 10 requests fail an error is returned.
func (c *Client) EnsureGrants(serviceName string) error {
	var response GrantsResponse
	doesNotExist := "id provided does not exist"
	isNotAService := "id provided is not a service"
	for i := 0; i < 10; i++ {
		err := c.Request("PUT", "/v1/grants/"+serviceName, nil, &response)
		if err == nil {
			if response.Error == doesNotExist || response.Error == isNotAService {
				return fmt.Errorf("Failed to create KMS grant for %s, got %s", serviceName, response.Error)
			} else if response.Grants.DecryptGrant && response.Grants.EncryptGrant {
				return nil
			}
		}
		log.Printf("Failed to create KMS grants for %s, trying again (err: %e, %s)", serviceName, err, response.Error)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("Failed to create KMS grants for %s", serviceName)
}
