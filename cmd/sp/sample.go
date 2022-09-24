package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewSampleCommand() *SampleCommand {
	c := &SampleCommand{
		fs: flag.NewFlagSet("sample", flag.ContinueOnError),
	}

	addUnitFlags(c.fs, &c.uf)

	c.fs.IntVar(&c.limit, "limit", 0, "Sample limit per sensor")
	return c
}

type SampleCommand struct {
	fs *flag.FlagSet

	uf unitFlags

	limit int
}

func (c *SampleCommand) Name() string {
	return c.fs.Name()
}

func (c *SampleCommand) Description() string {
	return "Query for samples"
}

func (c *SampleCommand) Run(args []string) error {
	if err := c.fs.Parse(args[1:]); err != nil {
		return err
	}

	fmtU, err := newUnitsFormatter(&c.uf)
	if err != nil {
		return err
	}

	ctx := context.Background()
	sc := newClient(ctx)

	filter := sensorpush.SampleQueryFilter{}
	if c.limit != 0 {
		filter.Limit = &c.limit
	}

	ss, err := sc.Sample.Query(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Print(fmtSamples(fmtU, ss))

	return nil
}

func fmtSamples(fmtU *unitsFormatter, ss *sensorpush.Samples) string {
	var b strings.Builder

	fmtAttrVal(&b, "Last Time", fmtU.Time(ss.LastTime), 0)
	fmtAttrVal(&b, "Status", ss.Status.String(), 0)
	fmtAttrVal(&b, "Total Samples", fmtU.Int(ss.TotalSamples), 0)
	fmtAttrVal(&b, "Total Sensors", fmtU.Int(ss.TotalSensors), 0)
	fmtAttrVal(&b, "Truncated", fmtU.Bool(ss.Truncated), 0)

	return b.String()
}
