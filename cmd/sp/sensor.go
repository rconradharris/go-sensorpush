package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrSensorList = "%-20s %-8s %6s %12s %10s %10s"
)

func cmdSensor() {
	flag.Parse()
	switch flag.Arg(1) {
	case "ls":
		cmdSensorList()
	case "show":
		cmdSensorShow()
	}
}

func usageSensorShow() {
	fmt.Println("Usage: sp sensor show NAME_OR_ID")
}

func cmdSensorShow() {
	nameOrID := flag.Arg(2)

	if nameOrID == "" {
		usageSensorShow()
		os.Exit(1)
	}

	ctx := context.Background()
	sc := newClient(ctx)

	ss, err := sc.Sensor.List(ctx, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	s := findSensorByNameOrID(ss, nameOrID)
	if s == nil {
		log.Fatalf("unable to find a sensor matching: '%s'", nameOrID)
	}
	fmt.Print(fmtSensorShow(s))
}

func fmtSensorShow(s *sensorpush.Sensor) string {
	var b strings.Builder
	fmtAttrVal(&b, "Name", s.Name, 0)
	fmtAttrVal(&b, "Type", s.Type.String(), 0)
	fmtAttrVal(&b, "Active", fmtBool(s.Active), 0)
	fmtAttrVal(&b, "Battery(V)", fmtVoltage(s.BatteryVoltage), 0)
	fmtAttrVal(&b, "Signal(dB)", fmtSignalStrength(s.RSSI), 0)
	fmtAttrVal(&b, "DeviceID", s.DeviceID, 0)
	fmtAttrVal(&b, "ID", s.ID, 0)

	c := s.Calibration
	fmtAttrValHeading(&b, "Calibration", 0)
	fmtAttrVal(&b, "Humidity", fmtHumidity(c.Humidity), 1)
	fmtAttrVal(&b, "Temperature", fmtTemperature(c.Temperature), 1)

	fmtAttrValHeading(&b, "Alerts", 0)

	ah := s.Alerts.Humidity
	fmtAttrValHeading(&b, "Humidity", 1)
	fmtAttrVal(&b, "Enabled", fmtBool(ah.Enabled), 2)
	fmtAttrVal(&b, "Max", fmtHumidity(ah.Max), 2)
	fmtAttrVal(&b, "Min", fmtHumidity(ah.Min), 2)

	at := s.Alerts.Temperature
	fmtAttrValHeading(&b, "Temperature", 1)
	fmtAttrVal(&b, "Enabled", fmtBool(at.Enabled), 2)
	fmtAttrVal(&b, "Max", fmtTemperature(at.Max), 2)
	fmtAttrVal(&b, "Min", fmtTemperature(at.Min), 2)

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

func cmdSensorList() {
	ctx := context.Background()

	sc := newClient(ctx)

	ss, err := sc.Sensor.List(ctx, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(fmtSensorHeading())
	for _, s := range ss {
		fmt.Println(fmtSensorList(s))
	}
}

func fmtSensorHeading() string {
	return fmt.Sprintf(fmtStrSensorList,
		"Name",
		"Type",
		"Active",
		"Battery(V)",
		"Signal(dB)",
		"DeviceID",
	)
}

func fmtSensorList(s *sensorpush.Sensor) string {
	return fmt.Sprintf(fmtStrSensorList,
		s.Name,
		s.Type,
		fmtBool(s.Active),
		fmtVoltage(s.BatteryVoltage),
		fmtSignalStrength(s.RSSI),
		s.DeviceID,
	)
}
