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

// ResponseGetTime represents the response received
type ResponseGetTime struct {
	Cmd   string           `json:"cmd"`
	Code  int              `json:"code"`
	Value respGetTimeValue `json:"value"`
}

type respGetTimeValue struct {
	Time respGetTimeTime `json:"Time"`
}

type respGetTimeTime struct {
	Day    int `json:"day"`
	Month  int `json:"mon"`
	Year   int `json:"year"`
	Hour   int `json:"hour"`
	Minute int `json:"min"`
	Second int `json:"sec"`
}
