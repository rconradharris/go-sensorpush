package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	sp "github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrSample = "%-15s %12s %10s %25s"
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
	c.fs.StringVar(&c.startTime, "start", "", "Start time in ISO format (ex: \"2006-01-02T15:04:05Z07:00\")")
	c.fs.StringVar(&c.stopTime, "stop", "", "Stop time in ISO format (ex: \"2006-01-02T15:04:05Z07:00\")")
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

func parseMeasures(str string) (sp.MeasureMap, error) {
	if str == "" {
		return nil, nil
	}
	items := parseCommaDelim(str)
	mm := sp.MeasureMap{}
	for _, s := range items {
		m, err := sp.ParseMeasure(s)
		if err != nil {
			return nil, err
		}
		mm.Add(m)
	}
	return mm, nil
}

func parseSensorIDs(sm sp.SensorMap, str string) ([]sp.SensorID, error) {
	if str == "" {
		return nil, nil
	}
	items := parseCommaDelim(str)
	ss := make([]sp.SensorID, 0, len(items))
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

	mm, err := parseMeasures(c.measures)
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

	filter := sp.SampleQueryFilter{
		Active:   c.active,
		Measures: mm.Measures(),
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
		fmt.Print(fmtSamplesVerbose(fmtU, mm, sm, ss))
	} else {
		fmt.Print(fmtSamples(fmtU, sm, ss))
	}

	return nil
}

func fmtSamples(fmtU *unitsFormatter, sm sp.SensorMap, ss *sp.Samples) string {
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

func fmtSample(b *strings.Builder, fmtU *unitsFormatter, sensor *sp.Sensor, s *sp.Sample) {
	fmt.Fprintf(b, fmtStrSample,
		sensor.Name,
		fmtU.Temperature(s.Temperature),
		fmtU.Humidity(s.Humidity),
		fmtU.Time(s.Observed),
	)
}

func fmtSamplesVerbose(fmtU *unitsFormatter, mm sp.MeasureMap, sm sp.SensorMap, ss *sp.Samples) string {
	var b strings.Builder

	fmtAttrVal(&b, "Last Time", fmtU.Time(ss.LastTime), 0)
	fmtAttrVal(&b, "Status", ss.Status.String(), 0)
	fmtAttrVal(&b, "Total Samples", fmtU.Int(ss.TotalSamples), 0)
	fmtAttrVal(&b, "Total Sensors", fmtU.Int(ss.TotalSensors), 0)
	fmtAttrVal(&b, "Truncated", fmtU.Bool(ss.Truncated), 0)

	fmtAttrValHeading(&b, "Sensor Samples", 0)

	for _, sn := range sm.SensorsAlpha() {
		of := newObservationFilter(mm, sn)
		samples, ok := ss.Sensors[sn.ID]
		if !ok {
			continue
		}
		fmtAttrValHeading(&b, sn.Name, 1)
		for _, sample := range samples {
			fmtSampleVerbose(&b, fmtU, of, sample)
		}
	}

	return b.String()
}

type observationFilter struct {
	mm sp.MeasureMap
	sn *sp.Sensor
	fs sp.SensorFeatureSet
}

func newObservationFilter(mm sp.MeasureMap, sn *sp.Sensor) observationFilter {
	return observationFilter{
		mm: mm,
		sn: sn,
		fs: sn.Type.Features(),
	}
}

func (of observationFilter) show(f sp.SensorFeature, m sp.Measure, defined bool) bool {
	if !of.fs.Has(f) {
		return false // Sensor doesn't support this feature
	}
	if of.mm.Has(m) {
		return true // User specifically requested this field
	}
	return defined // Otherwise, does it have a sensible value?
}

func fmtSampleVerbose(b *strings.Builder, fmtU *unitsFormatter, of observationFilter, s *sp.Sample) {
	fmtAttrVal(b, "Observed", fmtU.Time(s.Observed), 2)
	if of.show(sp.SensorFeatureBarometricPressure, sp.MeasureAltitude, s.Altitude != nil) {
		fmtAttrVal(b, "Altitude", fmtU.Distance(s.Altitude), 3)
	}
	if of.show(sp.SensorFeatureBarometricPressure, sp.MeasureBarometricPressure, s.BarometricPressure != nil) {
		fmtAttrVal(b, "Baro Pressure", fmtU.BarometricPressure(s.BarometricPressure), 3)
	}
	if of.show(sp.SensorFeatureDewPoint, sp.MeasureDewPoint, s.DewPoint != nil) {
		fmtAttrVal(b, "Dew Point", fmtU.Temperature(s.DewPoint), 3)
	}
	if of.show(sp.SensorFeatureHumidity, sp.MeasureHumidity, s.Humidity != nil) {
		fmtAttrVal(b, "Humidity", fmtU.Humidity(s.Humidity), 3)
	}
	if of.show(sp.SensorFeatureTemperature, sp.MeasureTemperature, s.Temperature != nil) {
		fmtAttrVal(b, "Temperature", fmtU.Temperature(s.Temperature), 3)
	}
	if of.show(sp.SensorFeatureVPD, sp.MeasureVPD, s.VPD != nil) {
		fmtAttrVal(b, "VPD", fmtU.VPD(s.VPD), 3)
	}
}
