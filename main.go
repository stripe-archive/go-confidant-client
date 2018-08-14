package main

import (
	"fmt"
	"net/http"

	"github.com/stripe/go-confidant-client/confidant"
	"github.com/stripe/go-confidant-client/kmsauth"
)

func main() {
	cluster := "cluster"
	authkey := "key"
	from := "username"
	to := "confidant"
	userType := "user"
	region := "region"
	url := "confidant-url"
	httpClient := &http.Client{}
	generator := kmsauth.NewTokenGenerator(authkey, to, from, userType, region)

	client := confidant.NewClient(url, httpClient, &generator)

	service, err := client.GetService("service-name")
	if err != nil {
		fmt.Printf("Error! %+v", err)
	}
	fmt.Printf("%+v", service)
}
