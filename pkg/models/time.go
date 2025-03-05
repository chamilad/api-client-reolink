// Package models contains the models that can be used to communciate with the camera API
package models

// RequestGetTime represents the GetTime request payload
type RequestGetTime struct {
	Cmd    string `json:"cmd"`
	Action int    `json:"action"`
}

// NewGetTimeRequest creates a GetTime payload
func NewGetTimeRequest() *[]RequestGetTime {
	return &[]RequestGetTime{
		{
			Cmd:    "GetTime",
			Action: 1,
		},
	}
}
