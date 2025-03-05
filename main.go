package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chamilad/api-client-reolink/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	host     = flag.String("host", "", "the ip address or the hostname of the camera with scheme, ex: https://192.168.1.100")
	username = flag.String("username", "", "the username to authenticate with")
	password = flag.String("password", "", "the password to authenticate with")
)

// Login logs into the camera and returns the token to be used on subsequent requests
func Login(host, username, password string) (string, error) {
	loginRequest := models.NewLoginRequest(username, password)
	payload, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("error while building login payload: %v", err)
	}

	log.Debug().Msg(fmt.Sprintf("json payload: %s", payload))

	loginURL := fmt.Sprintf("%s/cgi-bin/api.cgi", host)
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error while building login request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("cmd", "Login")
	q.Add("token", "null")
	req.URL.RawQuery = q.Encode()

	// todo: flag base
	insecureTr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	client := &http.Client{
		Transport: insecureTr,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while calling login request: %v", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading login response: %v", err)
	}

	fmt.Println("login response body: ", string(respBody))
	return string(respBody), nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// get ip addr, username, and password with cli args
	flag.Parse()
	_, err := Login(strings.TrimSpace(*host), strings.TrimSpace(*username), strings.TrimSpace(*password))
	if err != nil {
		fmt.Printf("error while logging in: %v", err)
	}
}
