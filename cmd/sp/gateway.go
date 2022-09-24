package main

import (
	"flag"
)

func NewGatewayCommand() *GatewayCommand {
	sc := &GatewayCommand{
		fs: flag.NewFlagSet("gateway", flag.ContinueOnError),
	}
	return sc
}

type GatewayCommand struct {
	fs *flag.FlagSet

	name string
}

func (c *GatewayCommand) Name() string {
	return c.fs.Name()
}

func (c *GatewayCommand) Description() string {
	return "Manage gateways"
}

func (c *GatewayCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	cmds := []Runner{
		NewGatewayListCommand(),
	}

	if err := DispatchCommand(cmds, args); err != nil {
		return err
	}
	return nil
}
