package confidant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestBody struct {
	Credentials      []string `json:"credentials"`
	BlindCredentials []string `json:"blind_credentials"`
	Account          string   `json:"account"`
	Enabled          bool     `json:"enabled"`
}

func (c *Client) Request(method string, path string, body *RequestBody, result interface{}) error {
	url := c.url + path

	if body != nil {
		// Marshal empty arrays instead of "null".  The Confidant API expects these to be arrays.
		if body.Credentials == nil {
			body.Credentials = make([]string, 0)
		}
		if body.BlindCredentials == nil {
			body.BlindCredentials = make([]string, 0)
		}
	}
	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	username := c.TokenGenerator.GetUsername()
	token, err := c.TokenGenerator.GetToken()
	if err != nil {
		return err
	}
	req.Header.Add("X-Auth-From", username)
	req.Header.Add("X-Auth-Token", token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("NotFound")
	} else if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("Forbidden")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Confidant Request Failed: got status code %v with body %s", resp.StatusCode, string(bodyBytes))
	}
	err = json.Unmarshal(bodyBytes, result)
	if err != nil {
		return err
	}
	return nil
}
