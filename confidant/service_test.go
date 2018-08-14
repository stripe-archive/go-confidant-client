package confidant

import (
	"reflect"
	"testing"
)

func testService(service *Service, expectedService *Service, t *testing.T) {
	if service.ID != expectedService.ID {
		t.Errorf("Incorrect service.ID: expected %s, got %s", expectedService.ID, service.ID)
	}
	if service.Enabled != expectedService.Enabled {
		t.Errorf("Incorrect service.Enabled: expected %t, got %t", expectedService.Enabled, service.Enabled)
	}
	if service.Account != expectedService.Account {
		t.Errorf("Incorrect service.Account: expected %s, got %s", expectedService.Account, service.Account)
	}
	if service.Revision != expectedService.Revision {
		t.Errorf("Incorrect service.Revision: expected %d, got %d", expectedService.Revision, service.Revision)
	}
	if !reflect.DeepEqual(service.BlindCredentials, expectedService.BlindCredentials) {
		t.Errorf("Blind credentials don't match: expected %v, got %v", expectedService.BlindCredentials, service.Credentials)
	}
	if !reflect.DeepEqual(service.Credentials, expectedService.Credentials) {
		t.Errorf("Credentials don't match: expected %v, got %v", expectedService.Credentials, service.Credentials)
	}
	if service.ModifiedBy != expectedService.ModifiedBy {
		t.Errorf("Incorrect service.ModifiedBy: expected %s, got %s", expectedService.ModifiedBy, service.ModifiedBy)
	}
	if service.ModifiedDate != expectedService.ModifiedDate {
		t.Errorf("Incorrect service.ModifiedDate: expected %s, got %s", expectedService.ModifiedDate, service.ModifiedDate)
	}
}

func TestGetServices(t *testing.T) {
	method := "GET"
	path := "/v1/services"
	service := Service{
		Enabled:          true,
		ID:               "service-name",
		Account:          "",
		BlindCredentials: nil,
		Credentials:      nil,
		Revision:         1,
		ModifiedBy:       "username",
		ModifiedDate:     "Wed, 25 Jul 2018 01:26:17 GMT",
	}
	expectedServices := Services{Services: []Service{service}}
	responses := map[string]interface{}{method + path: expectedServices}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	services, err := c.GetServices()
	if err != nil {
		t.Errorf("Could not get services: %e", err)
	}
	if len(services.Services) != len(expectedServices.Services) {
		t.Errorf("Incorrect number of services: expected %d, got %d", len(services.Services), len(expectedServices.Services))
	}
	service = services.Services[0]
	expectedService := expectedServices.Services[0]
	testService(&service, &expectedService, t)
}

func TestGetService(t *testing.T) {
	serviceName := "service-name"
	method := "GET"
	path := "/v1/services/" + serviceName
	credential := Credential{
		CredentialPairs: map[string]string{
			"key": "value",
		},
		Enabled:  true,
		ID:       "1",
		Name:     "name",
		Revision: 1,
	}
	credentials := []*Credential{&credential}
	expectedService := Service{
		ID:               serviceName,
		Enabled:          true,
		Account:          "",
		BlindCredentials: make([]*Credential, 0),
		Credentials:      credentials,
		Revision:         1,
		ModifiedBy:       "username",
		ModifiedDate:     "Wed, 25 Jul 2018 01:26:17 GMT",
	}
	responses := map[string]interface{}{method + path: expectedService}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	service, err := c.GetService(serviceName)
	if err != nil {
		t.Errorf("Could not get service: %e", err)
	}
	if _, ok := c.services[serviceName]; !ok {
		t.Errorf("client.services does not contain %s", serviceName)
	}
	testService(service, &expectedService, t)
}

func TestCreateService(t *testing.T) {
	serviceName := "foo"
	path := "/v1/services/" + serviceName
	credential := Credential{
		CredentialPairs: map[string]string{
			"key": "value",
		},
		Enabled:  true,
		ID:       "1",
		Name:     "name",
		Revision: 1,
	}
	credentials := []*Credential{&credential}
	expectedService := Service{
		ID:               serviceName,
		Enabled:          true,
		Account:          "",
		BlindCredentials: make([]*Credential, 0),
		Credentials:      credentials,
		Revision:         0,
		ModifiedBy:       "username",
		ModifiedDate:     "Wed, 25 Jul 2018 01:26:17 GMT",
	}
	responses := make(map[string]interface{})
	responses["PUT"+path] = ServiceResponse{Result: true, Service: expectedService}
	responses["GET/v1/services/"+serviceName] = ServiceResponse{Result: false, Error: "Service Doesn't Exist"}
	responses["GET/v1/roles"] = Roles{Roles: []string{serviceName}}
	responses["GET/v1/credentials"] = CredentialResponse{Credentials: []Credential{credential}}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	service, err := c.CreateService(serviceName, []string{"name"})
	if err != nil {
		t.Errorf("Could not create service: %e", err)
	}
	if _, ok := c.services[serviceName]; !ok {
		t.Errorf("client.services does not contain %s", serviceName)
	}
	testService(service, &expectedService, t)
}

func TestSetServiceCredentials(t *testing.T) {
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
		ID:           serviceName,
		Credentials:  []*Credential{&initialCredential},
		Revision:     1,
		ModifiedBy:   "username",
		ModifiedDate: "Wed, 25 Jul 2018 01:26:17 GMT",
	}
	expectedService := Service{
		ID:           serviceName,
		Credentials:  []*Credential{&newCredential},
		Revision:     2,
		ModifiedBy:   "username",
		ModifiedDate: "Wed, 25 Jul 2018 01:26:17 GMT",
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
	service, err := c.SetServiceCredentials(serviceName, []string{"new"})
	if err != nil {
		t.Errorf("Could not update service: %e", err)
	}
	if _, ok := c.services[serviceName]; !ok {
		t.Errorf("client.services does not contain %s", serviceName)
	}
	testService(service, &expectedService, t)
}

func TestUpdateServiceCredentials(t *testing.T) {
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
		ID:           serviceName,
		Credentials:  []*Credential{&initialCredential},
		Revision:     1,
		ModifiedBy:   "username",
		ModifiedDate: "Wed, 25 Jul 2018 01:26:17 GMT",
	}
	expectedService := Service{
		ID:           serviceName,
		Credentials:  []*Credential{&initialCredential, &newCredential},
		Revision:     2,
		ModifiedBy:   "username",
		ModifiedDate: "Wed, 25 Jul 2018 01:26:17 GMT",
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
	service, err := c.UpdateServiceCredentials(serviceName, []string{"new"}, []string{})
	if err != nil {
		t.Errorf("Could not update service: %e", err)
	}
	if _, ok := c.services[serviceName]; !ok {
		t.Errorf("client.services does not contain %s", serviceName)
	}
	testService(service, &expectedService, t)
}

func TestEnableService(t *testing.T) {
	serviceName := "foo"
	path := "/v1/services/" + serviceName
	initialCredential := Credential{
		ID:   "1",
		Name: "initial",
	}
	initialService := Service{
		Enabled:     false,
		Credentials: []*Credential{&initialCredential},
	}
	expectedService := Service{
		Enabled:     true,
		Credentials: []*Credential{&initialCredential},
	}
	grants := Grants{EncryptGrant: true, DecryptGrant: true}
	responses := make(map[string]interface{})
	responses["PUT"+path] = expectedService
	responses["GET/v1/services/"+serviceName] = ServiceResponse{Result: true, Service: initialService}
	responses["PUT/v1/grants/"+serviceName] = GrantsResponse{Grants: grants}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	service, err := c.EnableService(serviceName)
	if err != nil {
		t.Errorf("Could not enable service: %e", err)
	}
	if _, ok := c.services[serviceName]; !ok {
		t.Errorf("client.services does not contain %s", serviceName)
	}
	testService(service, &expectedService, t)
}

func TestDisableService(t *testing.T) {
	serviceName := "foo"
	path := "/v1/services/" + serviceName
	initialCredential := Credential{
		ID:   "1",
		Name: "initial",
	}
	initialService := Service{
		Enabled:     true,
		Credentials: []*Credential{&initialCredential},
	}
	expectedService := Service{
		Enabled:     false,
		Credentials: []*Credential{},
	}
	responses := make(map[string]interface{})
	responses["PUT"+path] = expectedService
	responses["GET/v1/services/"+serviceName] = ServiceResponse{Result: true, Service: initialService}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	service, err := c.DisableService(serviceName)
	if err != nil {
		t.Errorf("Could not disable service: %e", err)
	}
	if _, ok := c.services[serviceName]; !ok {
		t.Errorf("client.services does not contain %s", serviceName)
	}
	testService(service, &expectedService, t)
}
