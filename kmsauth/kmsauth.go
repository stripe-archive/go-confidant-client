package kmsauth

import (
	"fmt"
	"time"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

type TokenGenerator struct {
	KeyID	string
	Context map[string]*string
	KMSClient kmsiface.KMSAPI
}

type Payload struct {
	NotBefore string `json:"not_before"`
	NotAfter string `json:"not_after"`
}

func NewTokenGenerator(keyID, to string, from string, userType string, region string) TokenGenerator {
	context := map[string]*string{
		"from": aws.String(from),
		"to":  aws.String(to),
		"user_type": aws.String(userType),
	}
	config := &aws.Config{Region: aws.String(region)}
	client := kms.New(session.New(), config)
	return TokenGenerator{
		KeyID: keyID,
		Context: context,
		KMSClient: client,
	}
}

func (g *TokenGenerator) GetUsername() string {
	userType := aws.StringValue(g.Context["user_type"])
	from := aws.StringValue(g.Context["from"])
	return fmt.Sprintf("%d/%s/%s", 2, userType, from)
}

func (g *TokenGenerator) GetToken() (string, error) {
	now := time.Now().UTC()
	format := "20060102T150405Z"
	start := now.Format(format)
	end := now.Add(time.Minute * 60).Format(format)
	plaintext, err := json.Marshal(Payload{
		NotBefore: start,
		NotAfter: end,
	})
	if err != nil {
		return "", err
	}
	encrypted, err := g.Encrypt(plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (g *TokenGenerator) Encrypt(plaintext []byte) ([]byte, error) {
	input := &kms.EncryptInput{
		Plaintext:         plaintext,
		EncryptionContext: g.Context,
		GrantTokens:       []*string{},
		KeyId:             aws.String(g.KeyID),
	}
	resp, err := g.KMSClient.Encrypt(input)
	if err != nil {
		return []byte(""), err
	}
	return resp.CiphertextBlob, nil
}
