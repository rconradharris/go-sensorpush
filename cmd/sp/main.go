package main

import (
	"context"
	"log"
	"os"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func newClient(ctx context.Context) *sensorpush.Client {
	sc := sensorpush.NewClient(nil)

	email := os.Getenv("SENSORPUSH_EMAIL")
	password := os.Getenv("SENSORPUSH_PASSWORD")

	if email == "" {
		log.Fatal("SENSORPUSH_EMAIL required")
	}

	if password == "" {
		log.Fatal("SENSORPUSH_PASSWORD required")
	}

	auth, err := sc.Auth.Authorize(ctx, email, password)
	if err != nil {
		log.Fatal(err.Error())
	}

	tok, err := sc.Auth.AccessToken(ctx, auth)
	if err != nil {
		log.Fatal(err.Error())
	}
	sc.UseAccessToken(tok)

	return sc
}

func main() {
	switch os.Args[1] {
	case "sensor":
		cmdSensor()
	}
}
