package main

import (
	"flag"
)

func NewSensorCommand() *SensorCommand {
	sc := &SensorCommand{
		fs: flag.NewFlagSet("sensor", flag.ContinueOnError),
	}
	return sc
}

type SensorCommand struct {
	fs *flag.FlagSet

	name string
}

func (c *SensorCommand) Name() string {
	return c.fs.Name()
}

func (c *SensorCommand) Description() string {
	return "Manage sensors"
}

func (c *SensorCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	cmds := []Runner{
		NewSensorListCommand(),
		NewSensorShowCommand(),
	}

	if err := DispatchCommand(cmds, args); err != nil {
		return err
	}
	return nil
}
