package confidant

import (
	"errors"
)

type Roles struct {
	Result bool     `json:"result"`
	Roles  []string `json:"roles"`
}

func sliceContains(slice []string, element string) bool {
	for _, name := range slice {
		if name == element {
			return true
		}
	}
	return false
}

// CheckRole checks if the service name is a valid IAM role.
// The roles are fetched with a GET request to /v1/roles.
func (c *Client) CheckRole(serviceName string) error {
	var roles Roles
	err := c.Request("GET", "/v1/roles", nil, &roles)
	if err != nil {
		return err
	}
	if !sliceContains(roles.Roles, serviceName) {
		return errors.New("Invalid IAM Role")
	}
	return nil
}
