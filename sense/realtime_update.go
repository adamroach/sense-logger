package sense

import (
	"fmt"
	"strings"
	"time"
)

type RealtimeUpdate struct {
	Payload RealtimeUpdatePayload `json:"payload"`
	Type    string                `json:"type"`
}

type RealtimeUpdatePayload struct {
	Voltage        []float64               `json:"voltage"`
	FrameNumber    int                     `json:"frame"`
	Devices        []Device                `json:"devices"`
	DefaultCost    int                     `json:"defaultCost"`
	Channels       []float64               `json:"channels"`
	FrequencyHz    float64                 `json:"hz"`
	TotalWatts     float64                 `json:"w"`
	C              int                     `json:"c"`
	GridWatts      int                     `json:"grid_w"`
	Stats          RealtimeUpdateStats     `json:"_stats"`
	PowerFlow      RealtimeUpdatePowerFlow `json:"power_flow"`
	DeviceWatts    int                     `json:"d_w"`
	EpochTimestamp int64                   `json:"epoch"`
}

type RealtimeUpdateStats struct {
	BytesReceived    float64 `json:"brcv"`
	MessagesReceived float64 `json:"mrcv"`
	MessagesSent     float64 `json:"msnd"`
}

type RealtimeUpdatePowerFlow struct {
	Grid []string `json:"grid"`
}

func (r *RealtimeUpdate) String() string {
	var sb strings.Builder
	idWidth := 0
	nameWidth := 0
	for _, device := range r.Payload.Devices {
		if len(device.ID) > idWidth {
			idWidth = len(device.ID)
		}
		if len(device.Name) > nameWidth {
			nameWidth = len(device.Name)
		}
	}
	sb.WriteString(fmt.Sprintf("%-*s | %-*s | %9s\n", idWidth, "ID", nameWidth, "Name", "Watts"))
	sb.WriteString(fmt.Sprintf("%-*s-+-%-*s-+-%s\n", idWidth, strings.Repeat("-", idWidth), nameWidth, strings.Repeat("-", nameWidth), strings.Repeat("-", 20)))
	for _, device := range r.Payload.Devices {
		sb.WriteString(fmt.Sprintf("%-*s | %-*s | %9.4f %v\n", idWidth, device.ID, nameWidth, device.Name, *device.Watts, device.Attrs))
		/*
			if device.Tags.MergedDevices != nil {
				ids := strings.SplitSeq(*device.Tags.MergedDevices, ",")
				for id := range ids {
					mergedDevice, err := client.GetDeviceByID(strings.TrimSpace(id))
					if err != nil {
						sb.WriteString(fmt.Sprintf("Error getting merged device: %v\n", err))
						continue
					}
					sb.WriteString(fmt.Sprintf("%-*s | ┗━▶ %-*s |\n", idWidth, "", nameWidth-4, mergedDevice.Name))
				}
			}
		*/
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Voltage:        %v\n", r.Payload.Voltage))
	sb.WriteString(fmt.Sprintf("Channels:       %v\n", r.Payload.Channels))
	//sb.WriteString(fmt.Sprintf("Frame Number:   %v\n", r.Payload.FrameNumber))
	//sb.WriteString(fmt.Sprintf("Default Cost:   %v\n", r.Payload.DefaultCost))
	sb.WriteString(fmt.Sprintf("Total Watts:    %v\n", r.Payload.TotalWatts))
	sb.WriteString(fmt.Sprintf("Device Watts:   %v\n", r.Payload.DeviceWatts))
	sb.WriteString(fmt.Sprintf("Grid Watts:     %v\n", r.Payload.GridWatts))
	sb.WriteString(fmt.Sprintf("Epoch Timestamp:%v\n", time.Unix(r.Payload.EpochTimestamp, 0)))
	sb.WriteString(fmt.Sprintf("Frequency Hz:   %v\n", r.Payload.FrequencyHz))
	sb.WriteString(fmt.Sprintf("C:              %v\n", r.Payload.C))
	//sb.WriteString(fmt.Sprintf("Power Flow:     %v\n", r.Payload.PowerFlow.Grid))
	return sb.String()
}
