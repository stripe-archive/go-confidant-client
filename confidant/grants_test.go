package confidant

import (
	"testing"
)

func TestGetGrants(t *testing.T) {
	serviceName := "service-name"
	method := "GET"
	path := "/v1/grants/" + serviceName
	expected := Grants{
		EncryptGrant: true,
		DecryptGrant: true,
	}
	response := GrantsResponse{
		Grants: expected,
	}
	responses := map[string]interface{}{method + path: response}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	grants, err := c.GetGrants(serviceName)
	if err != nil {
		t.Errorf("Could not get grants for service %s: %e", serviceName, err)
	}
	if *grants != expected {
		t.Errorf("Grants don't match: expected %v, got %v", expected, grants)
	}
}

func TestEnsureGrants(t *testing.T) {
	serviceName := "service-name"
	method := "PUT"
	path := "/v1/grants/" + serviceName
	expected := Grants{
		EncryptGrant: true,
		DecryptGrant: true,
	}
	response := GrantsResponse{
		Grants: expected,
	}
	responses := map[string]interface{}{method + path: response}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	err := c.EnsureGrants(serviceName)
	if err != nil {
		t.Errorf("Could not ensure grants for service %s: %e", serviceName, err)
	}
}
