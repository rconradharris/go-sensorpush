package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
// fmtStrSensorList = "%-20s %-8s %6s %12s %10s %10s"
)

func NewGatewayListCommand() *GatewayListCommand {
	c := &GatewayListCommand{
		fs: flag.NewFlagSet("ls", flag.ContinueOnError),
	}
	return c
}

type GatewayListCommand struct {
	fs *flag.FlagSet
}

func (c *GatewayListCommand) Name() string {
	return c.fs.Name()
}

func (c *GatewayListCommand) Description() string {
	return "List gateways"
}

func (c *GatewayListCommand) Run(args []string) error {
	if err := c.fs.Parse(args[1:]); err != nil {
		return err
	}

	ctx := context.Background()

	sc := newClient(ctx)

	gs, err := sc.Gateway.List(ctx)
	if err != nil {
		return err
	}

	//fmt.Println(fmtSensorHeading())

	for _, g := range gs {
		fmt.Println(fmtGatewayList(g))
	}

	return nil
}

/*
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
*/

func fmtGatewayList(g *sensorpush.Gateway) string {
	return fmt.Sprintf("%s %s", g.Name, g.ID)
	/*
		return fmt.Sprintf(fmtStrSensorList,
			s.Name,
			s.Type,
			fmtU.Bool(s.Active),
			fmtU.Voltage(s.BatteryVoltage),
			fmtU.SignalStrength(s.RSSI),
			s.DeviceID,
		)
	*/
}
