package confidant

import (
	"fmt"
	"net/http"
	"os"

	"github.com/stripe/go-confidant-client/kmsauth"
)

func initClient() *Client {
	authkey := "key"
	to := "ConfidantServer"
	from := "username"
	userType := "user"
	region := "us-east-1"
	proxy := fmt.Sprintf("%s/.proxy", os.Getenv("HOME"))
	url := "confidant-url"
	httpClient := &http.Client{
		Transport: UnixProxy(proxy),
	}
	generator := kmsauth.NewTokenGenerator(authkey, to, from, userType, region)
	c := NewClient(url, httpClient, &generator)
	return &c
}
