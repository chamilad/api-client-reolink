package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chamilad/api-client-reolink/pkg/models"
	"github.com/rs/zerolog/log"
)

// RLClient is the type to use to interact with the camera
type RLClient struct {
	host   string
	token  string
	client *http.Client
}

// NewClient builds a new RLClient struct
func NewClient(host string, insecure bool) (*RLClient, error) {
	// todo: validate host

	insecureTr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		}}

	client := &http.Client{
		Transport: insecureTr,
	}

	return &RLClient{
		host:   strings.TrimRight(host, "/"),
		client: client,
	}, nil
}

// NewAPIRequest builds an http.Request struct based on the cmd and the payload
func (c *RLClient) NewAPIRequest(cmd string, payload *bytes.Buffer) (*http.Request, error) {
	url := fmt.Sprintf("%s/cgi-bin/api.cgi", c.host)
	// todo: check payload = nil case
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("error while building request: %v", err)
	}

	// payload is always app/json
	req.Header.Add("Content-Type", "application/json")

	// cmd parameter also needs to be setup
	q := req.URL.Query()
	q.Add("cmd", cmd)

	// the Login command needs no token
	if cmd != "Login" {
		log.Debug().Msg(fmt.Sprintf("adding token: %s", c.token))
		q.Add("token", c.token)
	}

	req.URL.RawQuery = q.Encode()

	return req, nil
}

// Login logs into the camera API
func (c *RLClient) Login(username, password string) (string, error) {
	loginRequest := models.NewLoginRequest(username, password)
	payload, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("error while building login payload: %v", err)
	}

	// loginURL := fmt.Sprintf("%s/cgi-bin/api.cgi", c.host)
	// req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(payload))
	// if err != nil {
	// 	return "", fmt.Errorf("error while building login request: %v", err)
	// }

	// req.Header.Add("Content-Type", "application/json")

	// q := req.URL.Query()
	// q.Add("cmd", "Login")
	// q.Add("token", "null")
	// req.URL.RawQuery = q.Encode()

	req, err := c.NewAPIRequest("Login", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error while building login request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while calling login request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading login response: %v", err)
	}

	// fmt.Println("login response body: ", string(respBody))
	loginResponse := []models.ResponseLogin{}
	err = json.Unmarshal(respBody, &loginResponse)
	if err != nil {
		return "", fmt.Errorf("error while building login response: %v", err)
	}

	if len(loginResponse) != 1 {
		return "", fmt.Errorf("maformed response received from login: %s", string(respBody))
	}

	token := loginResponse[0].Value.Token.Name
	c.token = token

	return string(respBody), nil
}

// GetTime reads the current time off the camera
func (c *RLClient) GetTime() (string, error) {
	payload, err := json.Marshal(models.NewGetTimeRequest())
	if err != nil {
		return "", fmt.Errorf("error while building gettime payload: %v", err)
	}

	req, err := c.NewAPIRequest("GetTime", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error while building gettime request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while making gettime request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading login response: %v", err)
	}

	fmt.Println("gettime response body: ", string(respBody))
	return string(respBody), nil
}
