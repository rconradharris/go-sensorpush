package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewSampleCommand() *SampleCommand {
	c := &SampleCommand{
		fs: flag.NewFlagSet("sample", flag.ContinueOnError),
	}

	c.fs.IntVar(&c.limit, "limit", 0, "Sample limit per sensor")
	return c
}

type SampleCommand struct {
	fs *flag.FlagSet

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

	fmt.Printf("ss => %+v\n", ss)

	return nil
}
