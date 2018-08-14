package kmsauth

import (
	"bytes"
	"time"
	"testing"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

type mockKMSClient struct {
	kmsiface.KMSAPI
	Resp kms.EncryptOutput
}

func (m *mockKMSClient) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return &m.Resp, nil
}

func TestGetUsername(t *testing.T) {
	key := "alias/authnz-production"
	to := "confidant-production"
	from := "terraform-provider-confidant"
	userType := "user"
	region := "us-east-1"
	expected := "2/user/terraform-provider-confidant"
	generator := NewTokenGenerator(key, to, from, userType, region)
	username := generator.GetUsername()
	if username != expected {
		t.Errorf(
			"Username generated for (from: %s, user_type: %s) was incorrect, got: %s, want: %s.",
			from, userType, username, expected,
		)
	}
}

func TestEncrypt(t *testing.T) {
	key := "alias/authnz-production"
	to := "confidant-production"
	from := "terraform-provider-confidant"
	userType := "user"
	region := "us-east-1"
	generator := NewTokenGenerator(key, to, from, userType, region)
	now := time.Now().UTC()
	format := "20180703T170301Z"
	start := now.Format(format)
	end := now.Add(time.Minute * 60).Format(format)
	plaintext, err := json.Marshal(Payload{
		NotBefore: start,
		NotAfter: end,
	})
	if err != nil {
		t.Errorf("Could not generate plaintext: %e", err)
	}
	expected := kms.EncryptOutput{
		CiphertextBlob: []byte("ZW5jcnlwdGVk"),
	}
	generator.KMSClient = &mockKMSClient{
		Resp: expected,
	}
	ciphertext, err := generator.Encrypt(plaintext)
	if err != nil {
		t.Errorf("Could not encrypt input: %e", err)
	}
	if !bytes.Equal(expected.CiphertextBlob, ciphertext) {
		t.Errorf("Encryption failed: expected %s as Ciphertextblob, got %s", expected.CiphertextBlob, ciphertext)
	}
}
