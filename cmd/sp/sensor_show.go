package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

func NewSensorShowCommand() *SensorShowCommand {
	c := &SensorShowCommand{
		fs: flag.NewFlagSet("show", flag.ContinueOnError),
	}
	addUnitFlags(c.fs, &c.uf)
	return c
}

type SensorShowCommand struct {
	fs *flag.FlagSet

	uf unitFlags
}

func (c *SensorShowCommand) Name() string {
	return c.fs.Name()
}

func (c *SensorShowCommand) Description() string {
	return "Show details for a sensor"
}

func (c *SensorShowCommand) Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("sensor show NAME_OR_ID")
	}

	nameOrID := args[1]

	if err := c.fs.Parse(args[2:]); err != nil {
		return err
	}

	fmtU, err := newUnitsFormatter(&c.uf)
	if err != nil {
		return err
	}

	ctx := context.Background()
	sc := newClient(ctx)

	ss, err := sc.Sensor.List(ctx, true)
	if err != nil {
		return err
	}

	s := findSensorByNameOrID(ss, nameOrID)
	if s == nil {
		return fmt.Errorf("unable to find a sensor matching: '%s'", nameOrID)
	}

	fmt.Print(fmtSensorShow(fmtU, s))

	return nil
}

func fmtSensorShow(fmtU *unitsFormatter, s *sensorpush.Sensor) string {
	var b strings.Builder
	fmtAttrVal(&b, "Name", s.Name, 0)
	fmtAttrVal(&b, "Type", s.Type.String(), 0)
	fmtAttrVal(&b, "Active", fmtU.Bool(s.Active), 0)
	fmtAttrVal(&b, "Battery", fmtU.Voltage(s.BatteryVoltage), 0)
	fmtAttrVal(&b, "Signal", fmtU.SignalStrength(s.RSSI), 0)
	fmtAttrVal(&b, "DeviceID", s.DeviceID, 0)
	fmtAttrVal(&b, "ID", s.ID, 0)

	c := s.Calibration
	fmtAttrValHeading(&b, "Calibration", 0)
	fmtAttrVal(&b, "Humidity", fmtU.HumidityDelta(c.HumidityDelta), 1)
	fmtAttrVal(&b, "Temperature", fmtU.TemperatureDelta(c.TemperatureDelta), 1)

	fmtAttrValHeading(&b, "Alerts", 0)

	ah := s.Alerts.Humidity
	fmtAttrValHeading(&b, "Humidity", 1)
	fmtAttrVal(&b, "Enabled", fmtU.Bool(ah.Enabled), 2)
	fmtAttrVal(&b, "Max", fmtU.Humidity(ah.Max), 2)
	fmtAttrVal(&b, "Min", fmtU.Humidity(ah.Min), 2)

	at := s.Alerts.Temperature
	fmtAttrValHeading(&b, "Temperature", 1)
	fmtAttrVal(&b, "Enabled", fmtU.Bool(at.Enabled), 2)
	fmtAttrVal(&b, "Max", fmtU.Temperature(at.Max), 2)
	fmtAttrVal(&b, "Min", fmtU.Temperature(at.Min), 2)

	return b.String()
}

// Returns sensor that matches:
//
// 1. Long ID exact match
// 2. Short ID exact match
// 3. Case-insensitive name
//
// Returns nil if no match is found
func findSensorByNameOrID(ss sensorpush.SensorSlice, nameOrID string) *sensorpush.Sensor {
	lowerName := strings.ToLower(nameOrID)
	for _, s := range ss {
		if s.ID == nameOrID {
			return s
		}
		if s.DeviceID == nameOrID {
			return s
		}
		if strings.ToLower(s.Name) == lowerName {
			return s
		}
	}
	return nil
}
