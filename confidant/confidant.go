package confidant

import (
	"net/http"

	"github.com/stripe/go-confidant-client/kmsauth"
)

func NewClient(url string, httpClient *http.Client, tokenGenerator *kmsauth.TokenGenerator) Client {
	client := Client{
		HttpClient:     httpClient,
		TokenGenerator: tokenGenerator,
		services:       make(map[string]*Service),
		url:            url,
	}
	return client
}

type Client struct {
	HttpClient     *http.Client
	TokenGenerator *kmsauth.TokenGenerator
	services       map[string]*Service
	url            string
}
