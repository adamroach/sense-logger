package rrd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/adamroach/sense-logger/sense"
	"github.com/ziutek/rrd"
)

const (
	MainFile   = "monitor.rrd"
	DeviceFile = "device.json"
	SampleRate = 1
	Heartbeat  = 2 * SampleRate
	VMax       = 250
	IMax       = 200 // maximum for 200-amp service
	WMax       = VMax * IMax
)

type Writer struct {
	directory  string
	lastReport time.Time
}

func NewWriter(directory string) (*Writer, error) {
	info, err := os.Stat(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory does not exist: %s", directory)
		}
		return nil, fmt.Errorf("error checking directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", directory)
	}

	return &Writer{
		directory: directory,
	}, nil
}

func (w *Writer) UpdateDeviceNames(devices []sense.Device) error {
	names := make(map[string]any)
	filePath := fmt.Sprintf("%s/%s", w.directory, DeviceFile)

	// Read the existing JSON file
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&names); err != nil {
			return fmt.Errorf("error decoding JSON file %v: %w", filePath, err)
		}
	}

	// Update the names map with device IDs and names
	for _, device := range devices {
		names[device.ID] = device.Name
	}

	// Write the updated names map back to the JSON file
	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating JSON file %v: %w", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(names); err != nil {
		return fmt.Errorf("error encoding JSON file %v: %w", filePath, err)
	}

	return nil
}

func (w *Writer) Write(update *sense.RealtimeUpdate) error {
	mainFilePath := fmt.Sprintf("%s/%s", w.directory, MainFile)

	reportTime := time.Unix(update.Payload.EpochTimestamp, 0)
	if reportTime.Sub(w.lastReport) < time.Duration(SampleRate)*time.Second {
		return nil
	}

	// Ensure the main RRD file exists
	if _, err := os.Stat(mainFilePath); os.IsNotExist(err) {
		fmt.Printf("%s does not exist; creating\n", mainFilePath)
		c := rrd.NewCreator(mainFilePath, time.Now(), 1)
		setCreatorParameters(c)
		/*
			Voltage:        [122.5346309842862 122.60344311894005]
			Channels:       [1480.872373590828 0]
			Total Watts:    1480.872373590828
			Device Watts:   1481
			Grid Watts:     1436
			Epoch Timestamp:2025-05-17 20:20:44 -0500 CDT
			Frequency Hz:   59.993350982666016
			C:              20
		*/
		c.DS("v1", "GAUGE", Heartbeat, -VMax, VMax)
		c.DS("v2", "GAUGE", Heartbeat, -VMax, VMax)
		c.DS("w1", "GAUGE", Heartbeat, -WMax, WMax)
		c.DS("w2", "GAUGE", Heartbeat, -WMax, WMax)
		c.DS("wt", "GAUGE", Heartbeat, -WMax, WMax)
		c.DS("wd", "GAUGE", Heartbeat, -WMax, WMax)
		c.DS("wg", "GAUGE", Heartbeat, -WMax, WMax)
		c.DS("hz", "GAUGE", Heartbeat, 0, 120)
		err := c.Create(true)
		if err != nil {
			return fmt.Errorf("error creating RRD file %v: %w", mainFilePath, err)
		}
	}

	// Update the main RRD file
	u := rrd.NewUpdater(mainFilePath)
	err := u.Update(reportTime,
		update.Payload.Voltage[0],
		update.Payload.Voltage[1],
		update.Payload.Channels[0],
		update.Payload.Channels[1],
		update.Payload.TotalWatts,
		update.Payload.DeviceWatts,
		update.Payload.GridWatts,
		update.Payload.FrequencyHz,
	)
	if err != nil {
		return fmt.Errorf("error updating RRD file %v: %w", mainFilePath, err)
	}

	// Update the device RRD files
	for _, device := range update.Payload.Devices {
		deviceFilePath := fmt.Sprintf("%s/%s.rrd", w.directory, device.ID)
		if _, err := os.Stat(deviceFilePath); os.IsNotExist(err) {
			fmt.Printf("%s does not exist; creating\n", deviceFilePath)
			c := rrd.NewCreator(deviceFilePath, time.Now(), 1)
			setCreatorParameters(c)
			c.DS("w", "GAUGE", Heartbeat, -WMax, WMax)
			c.DS("i", "GAUGE", Heartbeat, 0, IMax)
			c.DS("v", "GAUGE", Heartbeat, -VMax, VMax)
			c.DS("e", "GAUGE", Heartbeat, 0, 1_000_000_000) // Don't know the actual range here
			c.DS("ao_w", "GAUGE", Heartbeat, -WMax, WMax)
			err := c.Create(true)
			if err != nil {
				return fmt.Errorf("error creating RRD file %v: %w", deviceFilePath, err)
			}
		}
		// Update the device RRD file
		u := rrd.NewUpdater(deviceFilePath)
		statusDetails := device.StatusDetails
		if statusDetails == nil {
			statusDetails = &sense.StatusDetails{}
		}
		watts := 0.0
		ao_watts := 0.0
		if device.Watts != nil {
			watts = *device.Watts
		}
		if device.AlwaysOnWatts != nil {
			ao_watts = *device.AlwaysOnWatts
		}
		err := u.Update(reportTime,
			watts,
			statusDetails.Current,
			statusDetails.Voltage,
			statusDetails.EnergyUsed,
			ao_watts,
		)
		if err != nil {
			return fmt.Errorf("error updating RRD file %v: %w", deviceFilePath, err)
		}
	}

	w.lastReport = reportTime
	return nil
}

func setCreatorParameters(c *rrd.Creator) {
	// Every second for a week = 604800 seconds
	c.RRA("AVERAGE", 0.9, 1, 604800)
	// Every minute for a year = 525600 minutes
	c.RRA("AVERAGE", 0.5, 60, 525600)
	// Every hour for 30 years = 262800 hours
	c.RRA("AVERAGE", 0.5, 3600, 262800)
	// Hourly high and low for 30 years
	c.RRA("MIN", 0.5, 3600, 262800)
	c.RRA("MAX", 0.5, 3600, 262800)
}
