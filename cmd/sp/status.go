package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewStatusCommand() *statusCommand {
	sc := &statusCommand{
		fs: flag.NewFlagSet("status", flag.ContinueOnError),
	}
	return sc
}

type statusCommand struct {
	fs *flag.FlagSet

	name string
}

func (c *statusCommand) Name() string {
	return c.fs.Name()
}

func (c *statusCommand) Description() string {
	return "Query for API status"
}

func (c *statusCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	cfg := unitsCfg{}
	fmtU := newUnitsFormatter(cfg)

	ctx := context.Background()
	sc := newClient(ctx)

	st, err := sc.Status.Get(ctx)
	if err != nil {
		return err
	}

	fmt.Print(fmtStatus(fmtU, st))

	return nil
}

func fmtStatus(fmtU *unitsFormatter, st *sensorpush.Status) string {
	var b strings.Builder

	fmtAttrVal(&b, "Deployed", fmtU.Time(st.Deployed), 0)
	fmtAttrVal(&b, "Message", st.Message, 0)
	fmtAttrVal(&b, "MS", fmt.Sprintf("%d", st.MS), 0)
	fmtAttrVal(&b, "Stack", st.Stack, 0)
	fmtAttrVal(&b, "Status", st.Status.String(), 0)
	fmtAttrVal(&b, "Time", fmtU.Time(st.Time), 0)
	fmtAttrVal(&b, "Version", st.Version, 0)

	return b.String()
}
