package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewGatewayShowCommand() *GatewayShowCommand {
	c := &GatewayShowCommand{
		fs: flag.NewFlagSet("show", flag.ContinueOnError),
	}
	return c
}

type GatewayShowCommand struct {
	fs *flag.FlagSet
}

func (c *GatewayShowCommand) Name() string {
	return c.fs.Name()
}

func (c *GatewayShowCommand) Description() string {
	return "Show details for a gateway"
}

func (c *GatewayShowCommand) Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("gateway show NAME_OR_ID")
	}

	nameOrID := args[1]

	if err := c.fs.Parse(args[2:]); err != nil {
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

	s := findGatewayByNameOrID(gs, nameOrID)
	if s == nil {
		return fmt.Errorf("unable to find a gateway matching: '%s'", nameOrID)
	}

	fmt.Print(fmtGatewayShow(fmtU, s))

	return nil
}

func fmtGatewayShow(fmtU *unitsFormatter, g *sensorpush.Gateway) string {
	var b strings.Builder

	fmtAttrVal(&b, "Name", g.Name, 0)
	fmtAttrVal(&b, "ID", g.ID, 0)
	fmtAttrVal(&b, "Last Alert", fmtU.Time(g.LastAlert), 0)
	fmtAttrVal(&b, "Last Seen", fmtU.Time(g.LastSeen), 0)
	fmtAttrVal(&b, "Message", g.Message, 0)
	fmtAttrVal(&b, "Paired", fmtU.Bool(g.Paired), 0)
	fmtAttrVal(&b, "Version", g.Version, 0)

	return b.String()
}

// Returns gateway that matches:
//
// 1. ID exact match
// 2. Case-insensitive name
//
// Returns nil if no match is found
func findGatewayByNameOrID(gs sensorpush.GatewaySlice, nameOrID string) *sensorpush.Gateway {
	lowerName := strings.ToLower(nameOrID)
	for _, g := range gs {
		if g.ID == nameOrID {
			return g
		}
		if strings.ToLower(g.Name) == lowerName {
			return g
		}
	}
	return nil
}
