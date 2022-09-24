package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewSampleCommand() *SampleCommand {
	c := &SampleCommand{
		fs: flag.NewFlagSet("sample", flag.ContinueOnError),
	}

	addUnitFlags(c.fs, &c.uf)

	c.fs.BoolVar(&c.active, "active", true, "Filter by active devices")
	c.fs.IntVar(&c.limit, "limit", 0, "Sample limit per sensor")
	c.fs.StringVar(&c.measures, "measures", "default",
		"Measures to include (\"alt\", \"baro\", \"default\", \"dew\", \"hum\", \"temp\", \"vpd\")")
	c.fs.StringVar(&c.startTime, "start", "", "Start time (ex: \"2006-01-02T15:04:05Z07:00\")")
	c.fs.StringVar(&c.stopTime, "stop", "", "Stop time (ex: \"2006-01-02T15:04:05Z07:00\")")
	return c
}

type SampleCommand struct {
	fs *flag.FlagSet

	uf unitFlags

	active    bool
	limit     int
	measures  string
	startTime string
	stopTime  string
}

func (c *SampleCommand) Name() string {
	return c.fs.Name()
}

func (c *SampleCommand) Description() string {
	return "Query for samples"
}

func parseMeasures(str string) ([]sensorpush.Measure, error) {
	def := false
	ms := []sensorpush.Measure{}
	for _, s := range strings.Split(str, ",") {
		s = strings.TrimSpace(s)
		if s == "default" {
			def = true
			continue
		}

		m, err := sensorpush.ParseMeasure(s)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	if def && len(ms) > 0 {
		return nil, fmt.Errorf("'default' cannot be used with other measures specified")
	}

	if def {
		return []sensorpush.Measure{}, nil
	}

	return ms, nil
}

func (c *SampleCommand) Run(args []string) error {
	if err := c.fs.Parse(args[1:]); err != nil {
		return err
	}

	measures, err := parseMeasures(c.measures)
	if err != nil {
		return err
	}

	fmtU, err := newUnitsFormatter(&c.uf)
	if err != nil {
		return err
	}

	ctx := context.Background()
	sc := newClient(ctx)

	filter := sensorpush.SampleQueryFilter{
		Active:   c.active,
		Measures: measures,
	}

	if c.limit != 0 {
		filter.Limit = &c.limit
	}

	if c.startTime != "" {
		t, err := time.Parse(time.RFC3339, c.startTime)
		if err != nil {
			return err
		}
		filter.StartTime = t
	}

	if c.stopTime != "" {
		t, err := time.Parse(time.RFC3339, c.stopTime)
		if err != nil {
			return err
		}
		filter.StopTime = t
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

	fmtAttrValHeading(&b, "Sensor Samples", 0)
	for _, id := range ss.Sensors.IDs() {
		samples := ss.Sensors[id]
		name := fmt.Sprintf("Sensor %s", id)
		fmtAttrValHeading(&b, name, 1)
		for _, s := range samples {
			fmtSample(&b, fmtU, s)
		}
	}

	return b.String()
}

func fmtSample(b *strings.Builder, fmtU *unitsFormatter, s *sensorpush.Sample) {
	fmtAttrVal(b, "Observed", fmtU.Time(s.Observed), 2)
	fmtAttrVal(b, "Altitude", fmtU.Distance(s.Altitude), 3)
	fmtAttrVal(b, "Baro Pressure", fmtU.BarometricPressure(s.BarometricPressure), 3)
	fmtAttrVal(b, "Dew Point", fmtU.Temperature(s.DewPoint), 3)
	fmtAttrVal(b, "Humidity", fmtU.Humidity(s.Humidity), 3)
	fmtAttrVal(b, "Temperature", fmtU.Temperature(s.Temperature), 3)
	fmtAttrVal(b, "VPD", fmtU.VPD(s.VPD), 3)
}
