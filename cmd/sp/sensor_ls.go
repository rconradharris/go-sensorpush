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
	c := &SensorListCommand{
		fs: flag.NewFlagSet("ls", flag.ContinueOnError),
	}

	c.fs.BoolVar(&c.active, "active", true, "Filter by active")
	return c
}

type SensorListCommand struct {
	fs *flag.FlagSet

	active bool
}

func (c *SensorListCommand) Name() string {
	return c.fs.Name()
}

func (c *SensorListCommand) Description() string {
	return "List sensors"
}

func (c *SensorListCommand) Run(args []string) error {
	if err := c.fs.Parse(args[1:]); err != nil {
		return err
	}

	ctx := context.Background()

	sc := newClient(ctx)

	ss, err := sc.Sensor.List(ctx, c.active)
	if err != nil {
		return err
	}

	cfg := unitsCfg{}
	fmtU := newUnitsFormatter(cfg)

	fmt.Println(fmtSensorHeading())

	for _, s := range ss {
		fmt.Println(fmtSensorList(fmtU, s))
	}

	return nil
}

func fmtSensorHeading() string {
	return fmt.Sprintf(fmtStrSensorList,
		"Name",
		"Type",
		"Active",
		"Battery",
		"Signal",
		"DeviceID",
	)
}

func fmtSensorList(fmtU *unitsFormatter, s *sensorpush.Sensor) string {
	return fmt.Sprintf(fmtStrSensorList,
		s.Name,
		s.Type,
		fmtBool(s.Active),
		fmtU.Voltage(s.BatteryVoltage),
		fmtU.SignalStrength(s.RSSI),
		s.DeviceID,
	)
}
