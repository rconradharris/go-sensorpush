package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

type Runner interface {
	Name() string
	Run(args []string) error
}

func DispatchCommand(cmds []Runner, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("You must pass a sub-command")
	}

	subcommand := args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			return cmd.Run(args[1:])
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

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
	cmds := []Runner{
		NewSampleCommand(),
		NewSensorCommand(),
		NewStatusCommand(),
	}

	if err := DispatchCommand(cmds, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
