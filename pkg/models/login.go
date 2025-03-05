// Package models contains the models that can be used to communciate with the camera API
package models

// RequestLogin represents the login payload to the Login action
type RequestLogin struct {
	Cmd    string      `json:"cmd"`
	Action int         `json:"action"`
	Param  interface{} `json:"param"`
}

// ParamLogin param block to be used for login requests
type ParamLogin struct {
	User User `json:"User"`
}

// User contains the actual credentials in the request
type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// NewLoginRequest builds a login request
func NewLoginRequest(username, password string) *[]RequestLogin {
	return &[]RequestLogin{
		{
			Cmd:    "Login",
			Action: 0,
			Param: ParamLogin{
				User: User{
					UserName: username,
					Password: password,
				},
			},
		},
	}
}

// ResponseLogin used to unmarshal the json response from Login command
type ResponseLogin struct {
	Cmd   string             `json:"cmd"`
	Code  int                `json:"code"`
	Value responseLoginValue `json:"value"`
}

type responseLoginValue struct {
	Token responseLoginToken `json:"Token"`
}

type responseLoginToken struct {
	LeaseTime int    `json:"leaseTime"`
	Name      string `json:"name"`
}
