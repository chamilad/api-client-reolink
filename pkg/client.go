package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
func (c *RLClient) GetTime() (time.Time, error) {
	payload, err := json.Marshal(models.NewGetTimeRequest())
	if err != nil {
		return time.Time{}, fmt.Errorf("error while building gettime payload: %v", err)
	}

	req, err := c.NewAPIRequest("GetTime", bytes.NewBuffer(payload))
	if err != nil {
		return time.Time{}, fmt.Errorf("error while building gettime request: %v", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while making gettime request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while reading gettime response: %v", err)
	}

	// todo: build a datetime representation here
	// fmt.Println("gettime response body: ", string(respBody))
	getTimeResp := []models.ResponseGetTime{}
	if err = json.Unmarshal(respBody, &getTimeResp); err != nil {
		return time.Time{}, fmt.Errorf("error while building gettime response: %v", err)
	}

	if len(getTimeResp) != 1 {
		return time.Time{}, fmt.Errorf("maformed response received from gettime: %s", string(respBody))
	}

	cameraTime := getTimeResp[0]

	// todo: needs this as an input
	// a kind of a buffer overflow can be seen for the time zone field in the returned value from the API
	// +13 is shown as -43200 with the range showing 43200-(-46800). My guess is +13 overflows
	// into being -43200.
	// This causes ambiguity between NZ DST and actual GMT-12 so the time zone has to be an input
	// This will probably force the tz db mount if run inside docker
	tempLocation, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		return time.Time{}, fmt.Errorf("error while trying to load time zone: %v", err)
	}

	return time.Date(
		cameraTime.Value.Time.Year,
		time.Month(cameraTime.Value.Time.Month),
		cameraTime.Value.Time.Day,
		cameraTime.Value.Time.Hour,
		cameraTime.Value.Time.Minute,
		cameraTime.Value.Time.Second,
		0,
		tempLocation), nil
}
