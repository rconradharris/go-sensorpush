package main

import (
	"context"
	"flag"
	"fmt"

	sp "github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrGatewayList = "%-20s %6s %-20s %s"
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

	fmtU, err := newUnitsFormatter(nil)
	if err != nil {
		return err
	}

	ctx := context.Background()

	sc := newClient(ctx)

	gs, err := sc.Gateway.List(ctx)
	if err != nil {
		return err
	}

	fmt.Println(fmtGatewayListHeading())

	for _, g := range gs {
		fmt.Println(fmtGatewayList(fmtU, g))
	}

	return nil
}

func fmtGatewayListHeading() string {
	return fmt.Sprintf(fmtStrGatewayList,
		"Name",
		"Paired",
		"Last Seen",
		"ID",
	)
}

func fmtGatewayList(fmtU *unitsFormatter, g *sp.Gateway) string {
	return fmt.Sprintf(fmtStrGatewayList,
		g.Name,
		fmtU.Bool(g.Paired),
		fmtU.Time(g.LastSeen),
		g.ID,
	)
}
