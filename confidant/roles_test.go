package confidant

import (
	"testing"
)

func TestSliceContains(t *testing.T) {
	str := "test"
	includes := []string{str}
	excludes := []string{}
	if !sliceContains(includes, str) {
		t.Errorf("Expected slice to contain %s, got %t", str, sliceContains(includes, str))
	}
	if sliceContains(excludes, str) {
		t.Errorf("Expected slice to not contain %s, got %t", str, sliceContains(excludes, str))
	}
}

func TestCheckRole(t *testing.T) {
	serviceName := "service-name"
	method := "GET"
	path := "/v1/roles"
	expected := Roles{
		Result: true,
		Roles:  []string{serviceName},
	}
	responses := map[string]interface{}{method + path: expected}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	err := c.CheckRole(serviceName)
	if err != nil {
		if err.Error() == "Invalid IAM Role" {
			t.Errorf("Expected roles to contain service %s", serviceName)
		} else {
			t.Errorf("Could not check role for service %s: %e", serviceName, err)
		}
	}
}
