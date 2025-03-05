package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chamilad/api-client-reolink/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	host     = flag.String("host", "", "the ip address or the hostname of the camera with scheme, ex: https://192.168.1.100")
	username = flag.String("username", "", "the username to authenticate with")
	password = flag.String("password", "", "the password to authenticate with")
	insecure = flag.Bool("insecure", true, "true = skip TLS verification of camera cert")
	// todo: prompt for password flag
	// todo: interactive repl-like loop
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	flag.Parse()

	client, _ := client.NewClient(strings.TrimSpace(*host), *insecure)
	_, err := client.Login(strings.TrimSpace(*username), strings.TrimSpace(*password))
	if err != nil {
		log.Error().Msg(fmt.Sprintf("error while logging in: %v", err))
		os.Exit(1)
	}

	client.GetTime()
}
