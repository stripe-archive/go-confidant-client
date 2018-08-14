package confidant

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/stripe/go-confidant-client/kmsauth"
)

type mockKMSClient struct {
	kmsiface.KMSAPI
	Resp kms.EncryptOutput
}

func (m *mockKMSClient) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return &m.Resp, nil
}

func createHandlerFunc(t *testing.T, expectedUsername string, expectedToken string, responses map[string]interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		key := r.Method + r.URL.EscapedPath()
		response, ok := responses[key]
		if !ok {
			t.Errorf("Response for %s %s not provided", r.Method, r.URL.EscapedPath())
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content type application/json, got %s", contentType)
		}
		username := r.Header.Get("X-Auth-From")
		token := r.Header.Get("X-Auth-Token")
		if expectedUsername != username {
			t.Errorf("Expected username %s, got %s", expectedUsername, username)
		}
		if expectedToken != token {
			t.Errorf("Expected token %s, got %s", expectedToken, token)
		}
		json.NewEncoder(w).Encode(response)
	}
}

func CreateMockClientAndServer(responses map[string]interface{}, t *testing.T) (*httptest.Server, *Client) {
	authkey := "key"
	from := "go-confidant-client"
	to := "confidant"
	userType := "user"
	region := "region"
	token := "token"
	httpClient := &http.Client{}
	generator := kmsauth.NewTokenGenerator(authkey, to, from, userType, region)
	expected := kms.EncryptOutput{
		CiphertextBlob: []byte(token),
	}
	generator.KMSClient = &mockKMSClient{
		Resp: expected,
	}
	encodedToken := base64.StdEncoding.EncodeToString([]byte(token))
	username := fmt.Sprintf("2/user/%s", from)
	handler := createHandlerFunc(t, username, encodedToken, responses)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	c := NewClient(ts.URL, httpClient, &generator)
	return ts, &c
}

func TestRequest(t *testing.T) {
	var services []Service
	method := "GET"
	path := "/v1/services"
	service := Service{
		Enabled:     true,
		ID:          "test",
		Revision:    1,
		Credentials: make([]*Credential, 0),
	}
	expected := []Service{service}
	responses := map[string]interface{}{method + path: expected}
	ts, c := CreateMockClientAndServer(responses, t)
	defer ts.Close()
	err := c.Request(method, path, nil, &services)
	if err != nil {
		t.Errorf("Got an error when making the Confidant request: %e", err)
	}
	if len(services) != 1 {
		t.Errorf("Expected one service, got %d", len(services))
	}
	if !services[0].Enabled {
		t.Errorf("Service: expected enabled true, got %t", services[0].Enabled)
	}
	if services[0].ID != "test" {
		t.Errorf("Expected service id 'test', got %s", services[0].ID)
	}
	if services[0].Revision != 1 {
		t.Errorf("Expected service revision 1, got %d", services[0].Revision)
	}
}
