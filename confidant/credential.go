package confidant

import "fmt"

type Credential struct {
	CredentialPairs map[string]string `json:"credential_pairs"`
	Enabled         bool              `json:"enabled"`
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Revision        int               `json:"revision"`
}

type CredentialResponse struct {
	Result      bool         `json:"result"`
	Error       string       `json:"error"`
	Credentials []Credential `json:"credentials"`
}

// FindCredentialsByName returns a list of credentials for the names provided.
// It fetches all credentials with a GET request to /v1/credentials
// and filters them with the provided names.
// If any credentials are missing, returns an error containing their names instead.
func (c *Client) FindCredentialsByName(names []string) ([]*Credential, error) {
	var response CredentialResponse
	err := c.Request("GET", "/v1/credentials", nil, &response)
	if err != nil {
		return nil, err
	}
	credentialsMap := make(map[string]*Credential)
	for i, v := range response.Credentials {
		credentialsMap[v.Name] = &response.Credentials[i]
	}
	credentials := make([]*Credential, 0, len(names))
	missing := make([]string, 0, len(names))
	for _, v := range names {
		credential := credentialsMap[v]
		if credential == nil {
			missing = append(missing, v)
		} else {
			credentials = append(credentials, credential)
		}
	}
	if len(missing) != 0 {
		return nil, fmt.Errorf("The following credentials do not exist: %+v", missing)
	}
	return credentials, nil
}

// getCredentialIDs takes a map or slice of credentials and returns their IDs
func getCredentialIDs(credentials interface{}) []string {
	switch x := credentials.(type) {
	case map[string]*Credential:
		credentialIDs := make([]string, 0, len(x))
		for id, _ := range x {
			credentialIDs = append(credentialIDs, id)
		}
		return credentialIDs
	case []*Credential:
		credentialIDs := make([]string, 0, len(x))
		for _, credential := range x {
			credentialIDs = append(credentialIDs, credential.ID)
		}
		return credentialIDs
	}
	return []string{}
}

// createCredentialMap takes a slice of initial credentials
// and credentials to add and credentials to remove
// and returns a map of IDs to credentials
func createCredentialMap(init []*Credential, add []*Credential, remove []*Credential) map[string]*Credential {
	credentials := make(map[string]*Credential)
	for _, credential := range init {
		credentials[credential.ID] = credential
	}
	for _, credential := range add {
		credentials[credential.ID] = credential
	}
	for _, credential := range remove {
		delete(credentials, credential.ID)
	}
	return credentials
}

// AssignCredential assigns a credential to a service
func (c *Client) AssignCredential(serviceName, credentialName string) error {
	credentials := []string{credentialName}
	_, err := c.UpdateServiceCredentials(serviceName, credentials, nil)
	return err
}

// UnassignCredential removes a credential from a service
func (c *Client) UnassignCredential(serviceName, credentialName string) error {
	credentials := []string{credentialName}
	_, err := c.UpdateServiceCredentials(serviceName, nil, credentials)
	return err
}
