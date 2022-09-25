package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	sp "github.com/rconradharris/go-sensorpush/sensorpush"
)

type Runner interface {
	Name() string
	Description() string
	Run(args []string) error
}

func usage(cmds []Runner, args []string) error {
	var b strings.Builder

	// Handle sub-command (right now arbitrary nesting isn't supported)
	cmdName := filepath.Base(args[0])
	if cmdName != "sp" {
		cmdName = fmt.Sprintf("sp %s", cmdName)
	}

	fmt.Fprintf(&b, "Usage: %s COMMAND\n\n", cmdName)
	fmt.Fprintf(&b, "Commands:\n")
	for _, cmd := range cmds {
		fmt.Fprintf(&b, "  %-15s%s\n", cmd.Name(), cmd.Description())
	}
	return fmt.Errorf(b.String())
}

func DispatchCommand(cmds []Runner, args []string) error {
	if len(args) < 2 {
		return usage(cmds, args)
	}

	subcommand := args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			return cmd.Run(args[1:])
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func newClient(ctx context.Context) *sp.Client {
	sc := sp.NewClient(nil)

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
		NewGatewayCommand(),
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
