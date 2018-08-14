package kmsauth_test

import (
	"fmt"

	"github.com/stripe/go-confidant-client/kmsauth"
)

func Example() {
	// KMS key to use for authentication
	key := "key"
	// The service being authenticated
	to := "confidant"
	// The user for whom the token is being generated
	from := "username"
	userType := "user"
	region := "us-east-1"
	generator := kmsauth.NewTokenGenerator(key, to, from, userType, region)
	username := generator.GetUsername()
	token, err := generator.GetToken()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(username, token)
}
