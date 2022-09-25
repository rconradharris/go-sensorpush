package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrSample = "%-15s %12s %10s %20s"
)

func NewSampleCommand() *SampleCommand {
	c := &SampleCommand{
		fs: flag.NewFlagSet("sample", flag.ContinueOnError),
	}

	addUnitFlags(c.fs, &c.uf)

	c.fs.BoolVar(&c.active, "active", true, "Filter by active devices")
	c.fs.IntVar(&c.limit, "limit", 0, "Sample limit per sensor")
	c.fs.StringVar(&c.measures, "measures", "",
		"Measures to include (\"alt\", \"baro\", \"dew\", \"hum\", \"temp\", \"vpd\")")
	c.fs.StringVar(&c.sensors, "sensors", "", "Sensors to include (ID or name)")
	c.fs.StringVar(&c.startTime, "start", "", "Start time (ex: \"2006-01-02T15:04:05Z07:00\")")
	c.fs.StringVar(&c.stopTime, "stop", "", "Stop time (ex: \"2006-01-02T15:04:05Z07:00\")")
	c.fs.BoolVar(&c.verbose, "verbose", false, "Enable verbose mode")
	return c
}

type SampleCommand struct {
	fs *flag.FlagSet

	uf unitFlags

	active    bool
	limit     int
	measures  string
	sensors   string
	startTime string
	stopTime  string
	verbose   bool
}

func (c *SampleCommand) Name() string {
	return c.fs.Name()
}

func (c *SampleCommand) Description() string {
	return "Query for samples"
}

func parseCommaDelim(s string) []string {
	items := strings.Split(s, ",")
	n := len(items)
	for i := 0; i < n; i++ {
		items[0] = strings.TrimSpace(items[0])
	}
	return items
}

func parseMeasures(str string) ([]sensorpush.Measure, error) {
	if str == "" {
		return nil, nil
	}
	items := parseCommaDelim(str)
	ms := make([]sensorpush.Measure, 0, len(items))
	for _, s := range items {
		m, err := sensorpush.ParseMeasure(s)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

func parseSensorIDs(sm sensorpush.SensorMap, str string) ([]sensorpush.SensorID, error) {
	if str == "" {
		return nil, nil
	}
	items := parseCommaDelim(str)
	ss := make([]sensorpush.SensorID, 0, len(items))
	for _, nameOrID := range items {
		s := findSensorByNameOrID(sm, nameOrID)
		if s == nil {
			return nil, fmt.Errorf("no sensor named '%s' found", nameOrID)
		}
		ss = append(ss, s.ID)
	}
	return ss, nil
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

	sm, err := sc.Sensor.List(ctx, nil)
	if err != nil {
		return err
	}

	sensorIDs, err := parseSensorIDs(sm, c.sensors)
	if err != nil {
		return err
	}

	filter := sensorpush.SampleQueryFilter{
		Active:   c.active,
		Measures: measures,
		Sensors:  sensorIDs,
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

	if c.verbose {
		fmt.Print(fmtSamplesVerbose(fmtU, sm, ss))
	} else {
		fmt.Print(fmtSamples(fmtU, sm, ss))
	}

	return nil
}

func fmtSamples(fmtU *unitsFormatter, sm sensorpush.SensorMap, ss *sensorpush.Samples) string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s\n", fmtSampleHeading())
	for _, sensor := range sm.SensorsAlpha() {
		samples, ok := ss.Sensors[sensor.ID]
		if !ok {
			continue
		}
		for _, sample := range samples {
			fmtSample(&b, fmtU, sensor, sample)
			fmt.Fprintf(&b, "\n")
		}
	}

	return b.String()
}

func fmtSampleHeading() string {
	return fmt.Sprintf(fmtStrSample,
		"Name",
		"Temperature",
		"Humdity",
		"Observed",
	)
}

func fmtSample(b *strings.Builder, fmtU *unitsFormatter, sensor *sensorpush.Sensor, s *sensorpush.Sample) {
	fmt.Fprintf(b, fmtStrSample,
		sensor.Name,
		fmtU.Temperature(s.Temperature),
		fmtU.Humidity(s.Humidity),
		fmtU.HumanTime(s.Observed),
	)
}

func fmtSamplesVerbose(fmtU *unitsFormatter, sm sensorpush.SensorMap, ss *sensorpush.Samples) string {
	var b strings.Builder

	fmtAttrVal(&b, "Last Time", fmtU.Time(ss.LastTime), 0)
	fmtAttrVal(&b, "Status", ss.Status.String(), 0)
	fmtAttrVal(&b, "Total Samples", fmtU.Int(ss.TotalSamples), 0)
	fmtAttrVal(&b, "Total Sensors", fmtU.Int(ss.TotalSensors), 0)
	fmtAttrVal(&b, "Truncated", fmtU.Bool(ss.Truncated), 0)

	fmtAttrValHeading(&b, "Sensor Samples", 0)

	for _, s := range sm.SensorsAlpha() {
		samples, ok := ss.Sensors[s.ID]
		if !ok {
			continue
		}
		fmtAttrValHeading(&b, s.Name, 1)
		for _, s := range samples {
			fmtSampleVerbose(&b, fmtU, s)
		}
	}

	return b.String()
}

func fmtSampleVerbose(b *strings.Builder, fmtU *unitsFormatter, s *sensorpush.Sample) {
	fmtAttrVal(b, "Observed", fmtU.Time(s.Observed), 2)
	if s.Altitude != nil {
		fmtAttrVal(b, "Altitude", fmtU.Distance(s.Altitude), 3)
	}
	if s.BarometricPressure != nil {
		fmtAttrVal(b, "Baro Pressure", fmtU.BarometricPressure(s.BarometricPressure), 3)
	}
	if s.DewPoint != nil {
		fmtAttrVal(b, "Dew Point", fmtU.Temperature(s.DewPoint), 3)
	}
	if s.Humidity != nil {
		fmtAttrVal(b, "Humidity", fmtU.Humidity(s.Humidity), 3)
	}
	if s.Temperature != nil {
		fmtAttrVal(b, "Temperature", fmtU.Temperature(s.Temperature), 3)
	}
	if s.VPD != nil {
		fmtAttrVal(b, "VPD", fmtU.VPD(s.VPD), 3)
	}
}
