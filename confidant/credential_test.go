package confidant

import (
	"reflect"
	"testing"
)

func TestFindCredentialsByName(t *testing.T) {
	method := "GET"
	path := "/v1/credentials"
	included := Credential{
		ID:   "1",
		Name: "included",
	}
	excluded := Credential{
		ID:   "2",
		Name: "excluded",
	}
	response := CredentialResponse{
		Credentials: []Credential{included, excluded},
	}
	responses := map[string]interface{}{method + path: response}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	credentials, err := c.FindCredentialsByName([]string{"included"})
	expected := []*Credential{&included}
	if err != nil {
		t.Errorf("Could not find credentials by name: %e", err)
	}
	if len(expected) != len(credentials) {
		t.Errorf("Expected %d credentials, got %d", len(expected), len(credentials))
	}
	if expected[0].ID != credentials[0].ID {
		t.Errorf("Expected %+v credential, got %+v", expected[0], credentials[0])
	}
	credentials, err = c.FindCredentialsByName([]string{"non-existant"})
	if err == nil || err.Error() != "The following credentials do not exist: [non-existant]" {
		t.Errorf("Expected error (The following credentials do not exist: [non-existant]), got %e", err)
	}
}

func TestGetCredentialIDs(t *testing.T) {
	ID := "test"
	expected := [1]string{ID}
	credentials := map[string]*Credential{ID: nil}
	credentialIDs := getCredentialIDs(credentials)
	if expected[0] != credentialIDs[0] {
		t.Errorf("Map: IDs don't match, expected %v, got %v", expected, credentialIDs)
	}
	credentialSlice := []*Credential{&Credential{ID: ID}}
	credentialIDs = getCredentialIDs(credentialSlice)
	if expected[0] != credentialIDs[0] {
		t.Errorf("Slice: IDs don't match, expected %v, got %v", expected, credentialIDs)
	}
}

func TestCreateCredentialMap(t *testing.T) {
	initID := "test"
	addID := "add"
	removeID := "test"
	init := &Credential{ID: initID}
	add := &Credential{ID: addID}
	remove := &Credential{ID: removeID}
	expected := map[string]*Credential{
		addID: add,
	}
	credentials := createCredentialMap([]*Credential{init}, []*Credential{add}, []*Credential{remove})
	if !reflect.DeepEqual(credentials, expected) {
		t.Errorf("Maps don't match, expected %v, got %v", expected, credentials)
	}
}

func TestAssignCredential(t *testing.T) {
	serviceName := "foo"
	path := "/v1/services/" + serviceName
	initialCredential := Credential{
		ID:   "1",
		Name: "initial",
	}
	newCredential := Credential{
		ID:   "2",
		Name: "new",
	}
	initialService := Service{
		Account:     "",
		Credentials: []*Credential{&initialCredential},
	}
	expectedService := Service{
		Account:          "",
		BlindCredentials: make([]*Credential, 0),
		Credentials:      []*Credential{&initialCredential, &newCredential},
	}
	grants := Grants{EncryptGrant: true, DecryptGrant: true}
	responses := make(map[string]interface{})
	responses["PUT"+path] = expectedService
	responses["GET/v1/services/"+serviceName] = ServiceResponse{Result: true, Service: initialService}
	responses["PUT/v1/grants/"+serviceName] = GrantsResponse{Grants: grants}
	responses["GET/v1/roles"] = Roles{Roles: []string{serviceName}}
	responses["GET/v1/credentials"] = CredentialResponse{Credentials: []Credential{initialCredential, newCredential}}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	err := c.AssignCredential(serviceName, "new")
	if err != nil {
		t.Errorf("Could not create service: %e", err)
	}
	service := c.services[serviceName]
	if len(service.Credentials) != len(expectedService.Credentials) {
		t.Errorf("Incorrect number of credentials: expected %d, got %d", len(expectedService.Credentials), len(service.Credentials))
	}
	if !reflect.DeepEqual(service.Credentials, expectedService.Credentials) {
		t.Errorf("Credentials don't match: expected %v, got %v", expectedService.Credentials, service.Credentials)
	}
}

func TestUnassignCredential(t *testing.T) {
	serviceName := "foo"
	path := "/v1/services/" + serviceName
	initialCredential := Credential{
		ID:   "1",
		Name: "initial",
	}
	initialService := Service{
		Account:     "",
		Credentials: []*Credential{&initialCredential},
	}
	expectedService := Service{
		Account:          "",
		BlindCredentials: make([]*Credential, 0),
		Credentials:      []*Credential{},
	}
	grants := Grants{EncryptGrant: true, DecryptGrant: true}
	responses := make(map[string]interface{})
	responses["PUT"+path] = expectedService
	responses["GET/v1/services/"+serviceName] = ServiceResponse{Result: true, Service: initialService}
	responses["PUT/v1/grants/"+serviceName] = GrantsResponse{Grants: grants}
	responses["GET/v1/roles"] = Roles{Roles: []string{serviceName}}
	responses["GET/v1/credentials"] = CredentialResponse{Credentials: []Credential{initialCredential}}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	err := c.UnassignCredential(serviceName, "initial")
	if err != nil {
		t.Errorf("Could not create service: %e", err)
	}
	service := c.services[serviceName]
	if len(service.Credentials) != len(expectedService.Credentials) {
		t.Errorf("Incorrect number of credentials: expected %d, got %d", len(expectedService.Credentials), len(service.Credentials))
	}
	if !reflect.DeepEqual(service.Credentials, expectedService.Credentials) {
		t.Errorf("Credentials don't match: expected %v, got %v", expectedService.Credentials, service.Credentials)
	}
}
