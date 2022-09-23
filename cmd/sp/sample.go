package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewSampleCommand() *SampleCommand {
	sc := &SampleCommand{
		fs: flag.NewFlagSet("sample", flag.ContinueOnError),
	}
	return sc
}

type SampleCommand struct {
	fs *flag.FlagSet
}

func (c *SampleCommand) Name() string {
	return c.fs.Name()
}

func (c *SampleCommand) Run(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	ctx := context.Background()
	sc := newClient(ctx)

	filter := sensorpush.SampleQueryFilter{}

	//limit := 1
	//filter.Limit = &limit

	ss, err := sc.Sample.Query(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Printf("ss => %+v\n", ss)

	return nil
}
