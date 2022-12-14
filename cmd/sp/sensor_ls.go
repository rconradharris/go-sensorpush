package main

import (
	"context"
	"flag"
	"fmt"

	sp "github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrSensorList = "%-20s %-8s %6s %12s %10s %10s"
)

func NewSensorListCommand() *SensorListCommand {
	c := &SensorListCommand{
		fs: flag.NewFlagSet("ls", flag.ContinueOnError),
	}

	addUnitFlags(c.fs, &c.uf)

	c.fs.BoolVar(&c.active, "active", true, "Filter by active")
	return c
}

type SensorListCommand struct {
	fs *flag.FlagSet

	uf unitFlags

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

	f := &sp.SensorListFilter{}
	if !c.active {
		f.Active = &c.active
	}

	sm, err := sc.Sensor.List(ctx, f)
	if err != nil {
		return err
	}

	fmtU, err := newUnitsFormatter(&c.uf)
	if err != nil {
		return err
	}

	fmt.Println(fmtSensorHeading())

	for _, s := range sm.SensorsAlpha() {
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

func fmtSensorList(fmtU *unitsFormatter, s *sp.Sensor) string {
	return fmt.Sprintf(fmtStrSensorList,
		s.Name,
		s.Type,
		fmtU.Bool(s.Active),
		fmtU.Voltage(s.BatteryVoltage),
		fmtU.SignalStrength(s.RSSI),
		s.DeviceID,
	)
}
