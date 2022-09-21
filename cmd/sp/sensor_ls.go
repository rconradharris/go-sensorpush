package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrSensorList = "%-20s %-8s %6s %12s %10s %10s"
)

func NewSensorListCommand() *SensorListCommand {
	sc := &SensorListCommand{
		fs: flag.NewFlagSet("ls", flag.ContinueOnError),
	}
	return sc
}

type SensorListCommand struct {
	fs *flag.FlagSet

	name string
}

func (c *SensorListCommand) Name() string {
	return c.fs.Name()
}

func (c *SensorListCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()

	sc := newClient(ctx)

	ss, err := sc.Sensor.List(ctx, true)
	if err != nil {
		return err
	}

	fmt.Println(fmtSensorHeading())

	for _, s := range ss {
		fmt.Println(fmtSensorList(s))
	}

	return nil
}

func fmtSensorHeading() string {
	return fmt.Sprintf(fmtStrSensorList,
		"Name",
		"Type",
		"Active",
		"Battery(V)",
		"Signal(dB)",
		"DeviceID",
	)
}

func fmtSensorList(s *sensorpush.Sensor) string {
	return fmt.Sprintf(fmtStrSensorList,
		s.Name,
		s.Type,
		fmtBool(s.Active),
		fmtVoltage(s.BatteryVoltage),
		fmtSignalStrength(s.RSSI),
		s.DeviceID,
	)
}
