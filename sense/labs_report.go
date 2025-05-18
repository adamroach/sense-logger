package sense

import (
	"fmt"
	"strings"
	"time"
)

type LabsReport struct {
	FaultDetectionJSON struct {
		MotorStalls     map[string]any `json:"motorStalls"`
		CloggedDryer    map[string]any `json:"cloggedDryer"`
		SolarDisruption map[string]any `json:"solarDisruption"`
		PowerQuality    struct {
			DisplaySupport bool `json:"displaySupport"`
			Data           struct {
				StartDate          time.Time              `json:"start_date"`
				EndDate            time.Time              `json:"end_date"`
				Length             int                    `json:"length"`
				IncrementInSeconds int                    `json:"increment_in_seconds"`
				Channel0           *ChannelData           `json:"channel0"`
				Channel1           *ChannelData           `json:"channel1"`
				Channel2           *ChannelData           `json:"channel2"`
				Latest             map[string]VoltageData `json:"latest"`
				Min                VoltageSummary         `json:"min"`
				Max                VoltageSummary         `json:"max"`
				Average            VoltageSummary         `json:"average"`
			} `json:"data"`
			ReportRange int `json:"reportRange"`
			Compare     struct {
				Count       int    `json:"count"`
				Frequency   string `json:"frequency"`
				Percent     int    `json:"percent"`
				LowPercent  int    `json:"lowPercent"`
				LowRange    string `json:"lowRange"`
				MedPercent  int    `json:"medPercent"`
				MedRange    string `json:"medRange"`
				HighPercent int    `json:"highPercent"`
				HighRange   string `json:"highRange"`
			} `json:"compare"`
		} `json:"powerQuality"`
	} `json:"fault_detection_json"`
	PowerQualityRawCSV            string `json:"power_quality_raw_csv"`
	MotorStallDailyCountsSVG      string `json:"motor_stall_daily_counts_svg"`
	MotorStallRawCSV              string `json:"motor_stall_raw_csv"`
	MotorStallPowermeterSampleSVG string `json:"motor_stall_powermeter_sample_svg"`
}

type ChannelData struct {
	VMin []float64 `json:"vmin"`
	VMax []float64 `json:"vmax"`
}

type VoltageData struct {
	Date time.Time `json:"date"`
	V0   float64   `json:"v0"`
	V1   float64   `json:"v1"`
	V2   *float64  `json:"v2"`
}

type VoltageSummary struct {
	V0 *float64 `json:"v0"`
	V1 *float64 `json:"v1"`
	V2 *float64 `json:"v2"`
}

func (lr *LabsReport) Summary() string {
	var sb strings.Builder

	sb.WriteString("Labs Report Summary\n")
	sb.WriteString(strings.Repeat("=", 50) + "\n")

	// Fault Detection JSON
	sb.WriteString("Fault Detection:\n")
	sb.WriteString("  Motor Stalls:\n")
	motorStallsKeyWidth := 0
	for key := range lr.FaultDetectionJSON.MotorStalls {
		if len(key) > motorStallsKeyWidth {
			motorStallsKeyWidth = len(key)
		}
	}
	for key, value := range lr.FaultDetectionJSON.MotorStalls {
		sb.WriteString(fmt.Sprintf("    %*s: %v\n", motorStallsKeyWidth, key, value))
	}

	sb.WriteString("  Clogged Dryer:\n")
	cloggedDryerKeyWidth := 0
	for key := range lr.FaultDetectionJSON.CloggedDryer {
		if len(key) > cloggedDryerKeyWidth {
			cloggedDryerKeyWidth = len(key)
		}
	}
	for key, value := range lr.FaultDetectionJSON.CloggedDryer {
		sb.WriteString(fmt.Sprintf("    %*s: %v\n", cloggedDryerKeyWidth, key, value))
	}

	sb.WriteString("  Solar Disruption:\n")
	solarDisruptionKeyWidth := 0
	for key := range lr.FaultDetectionJSON.SolarDisruption {
		if len(key) > solarDisruptionKeyWidth {
			solarDisruptionKeyWidth = len(key)
		}
	}
	for key, value := range lr.FaultDetectionJSON.SolarDisruption {
		sb.WriteString(fmt.Sprintf("    %*s: %v\n", solarDisruptionKeyWidth, key, value))
	}

	// Power Quality
	sb.WriteString("  Power Quality:\n")
	sb.WriteString(fmt.Sprintf("    Display Support: %v\n", lr.FaultDetectionJSON.PowerQuality.DisplaySupport))
	sb.WriteString("    Data:\n")
	sb.WriteString(fmt.Sprintf("      Start Date: %v\n", lr.FaultDetectionJSON.PowerQuality.Data.StartDate))
	sb.WriteString(fmt.Sprintf("      End Date:   %v\n", lr.FaultDetectionJSON.PowerQuality.Data.EndDate))
	sb.WriteString(fmt.Sprintf("      Length:     %d\n", lr.FaultDetectionJSON.PowerQuality.Data.Length))
	sb.WriteString(fmt.Sprintf("      Increment In Seconds: %d\n", lr.FaultDetectionJSON.PowerQuality.Data.IncrementInSeconds))

	// Channel Data
	for i, channel := range []*ChannelData{
		lr.FaultDetectionJSON.PowerQuality.Data.Channel0,
		lr.FaultDetectionJSON.PowerQuality.Data.Channel1,
		lr.FaultDetectionJSON.PowerQuality.Data.Channel2,
	} {
		if channel != nil {
			sb.WriteString(fmt.Sprintf("      Channel %d:\n", i))
			sb.WriteString(fmt.Sprintf("        VMin: %v values\n", len(channel.VMin)))
			sb.WriteString(fmt.Sprintf("        VMax: %v valued\n", len(channel.VMax)))
		}
	}

	// Latest Voltage Data
	sb.WriteString("      Latest Voltage Data:\n")
	latestKeyWidth := 0
	for key := range lr.FaultDetectionJSON.PowerQuality.Data.Latest {
		if len(key) > latestKeyWidth {
			latestKeyWidth = len(key)
		}
	}
	for key, value := range lr.FaultDetectionJSON.PowerQuality.Data.Latest {
		sb.WriteString(fmt.Sprintf("        %*s: Date: %v, V0: %v, V1: %v, V2: %s\n",
			latestKeyWidth, key, value.Date, value.V0, value.V1, pointerToString(value.V2)))
	}

	// Voltage Summaries
	sb.WriteString("      Voltage Summaries:\n")
	sb.WriteString(fmt.Sprintf("        Min: V0: %v, V1: %v, V2: %v\n",
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Min.V0),
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Min.V1),
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Min.V2)))
	sb.WriteString(fmt.Sprintf("        Max: V0: %v, V1: %v, V2: %v\n",
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Max.V0),
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Max.V1),
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Max.V2)))
	sb.WriteString(fmt.Sprintf("        Average: V0: %v, V1: %v, V2: %v\n",
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Average.V0),
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Average.V1),
		pointerToString(lr.FaultDetectionJSON.PowerQuality.Data.Average.V2)))

	// Report Range and Compare
	sb.WriteString(fmt.Sprintf("    Report Range: %d\n", lr.FaultDetectionJSON.PowerQuality.ReportRange))
	sb.WriteString("    Compare:\n")
	sb.WriteString(fmt.Sprintf("      Count: %d\n", lr.FaultDetectionJSON.PowerQuality.Compare.Count))
	sb.WriteString(fmt.Sprintf("      Frequency: %s\n", lr.FaultDetectionJSON.PowerQuality.Compare.Frequency))
	sb.WriteString(fmt.Sprintf("      Percent: %d\n", lr.FaultDetectionJSON.PowerQuality.Compare.Percent))
	sb.WriteString(fmt.Sprintf("      Low Percent: %d, Low Range: %s\n", lr.FaultDetectionJSON.PowerQuality.Compare.LowPercent, lr.FaultDetectionJSON.PowerQuality.Compare.LowRange))
	sb.WriteString(fmt.Sprintf("      Med Percent: %d, Med Range: %s\n", lr.FaultDetectionJSON.PowerQuality.Compare.MedPercent, lr.FaultDetectionJSON.PowerQuality.Compare.MedRange))
	sb.WriteString(fmt.Sprintf("      High Percent: %d, High Range: %s\n", lr.FaultDetectionJSON.PowerQuality.Compare.HighPercent, lr.FaultDetectionJSON.PowerQuality.Compare.HighRange))

	// Other Fields
	/*
		sb.WriteString("Other Fields:\n")
		sb.WriteString(fmt.Sprintf("  Power Quality Raw CSV: %s\n", lr.PowerQualityRawCSV))
		sb.WriteString(fmt.Sprintf("  Motor Stall Daily Counts SVG: %s\n", lr.MotorStallDailyCountsSVG))
		sb.WriteString(fmt.Sprintf("  Motor Stall Raw CSV: %s\n", lr.MotorStallRawCSV))
		sb.WriteString(fmt.Sprintf("  Motor Stall Powermeter Sample SVG: %s\n", lr.MotorStallPowermeterSampleSVG))
	*/

	return sb.String()
}

func pointerToString(ptr *float64) string {
	if ptr == nil {
		return "-"
	}
	return fmt.Sprintf("%v", *ptr)
}
