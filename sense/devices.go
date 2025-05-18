package sense

import "time"

type Devices struct {
	Devices            []Device `json:"devices"`
	DeviceDataChecksum string   `json:"device_data_checksum"`
}

type Device struct {
	// Mandatory fields
	ID   string     `json:"id"`
	Name string     `json:"name"`
	Tags DeviceTags `json:"tags"`

	// Device endpoint fields
	GivenLocation   *string          `json:"given_location,omitempty"`
	GivenMake       *string          `json:"given_make,omitempty"`
	GivenModel      *string          `json:"given_model,omitempty"`
	Icon            *string          `json:"icon,omitempty"`
	LastMergedState *LastMergedState `json:"last_merged_state,omitempty"`
	LastState       *string          `json:"last_state,omitempty"`
	LastStateTime   *time.Time       `json:"last_state_time,omitempty"`
	Location        *string          `json:"location,omitempty"`
	Make            *string          `json:"make,omitempty"`
	Model           *string          `json:"model,omitempty"`
	MonitorID       *int             `json:"monitor_id,omitempty"`
	StandbyConfig   *StandbyConfig   `json:"standbyConfig,omitempty"`

	// Realtime update fields
	AlwaysOnState *bool          `json:"ao_st,omitempty"`
	AlwaysOnWatts *float64       `json:"ao_w,omitempty"`
	Attrs         []Attr         `json:"attrs,omitempty"`
	Cirumference  *int           `json:"c,omitempty"`
	StatusDetails *StatusDetails `json:"sd,omitempty"`
	Watts         *float64       `json:"w,omitempty"`
}

type Attr string

const (
	AttrIdle Attr = "Idle"
)

type StatusDetails struct {
	Watts      float64 `json:"w"`
	Current    float64 `json:"i"`
	Voltage    float64 `json:"v"`
	EnergyUsed float64 `json:"e"`
}

type LastMergedState struct {
	CombinedState string `json:"combinedState"`
	StateInts     []int  `json:"stateInts"`
}

type StandbyConfig struct {
	DefaultStandbyHysteresisSec int `json:"default_standby_hysteresis_sec"`
	DefaultStandbyThresholdWatt int `json:"default_standby_threshold_watt"`
	StandbyHysteresisSec        int `json:"standby_hysteresis_sec"`
	StandbyThresholdWatt        int `json:"standby_threshold_watt"`
}

type DeviceTags struct {
	Alertable                   *string    `json:"Alertable,omitempty"`
	AlwaysOn                    *string    `json:"AlwaysOn,omitempty"`
	ControlCapabilities         []string   `json:"ControlCapabilities,omitempty"`
	DCMActive                   *string    `json:"DCMActive,omitempty"`
	DUID                        *string    `json:"DUID,omitempty"`
	DateCreated                 *time.Time `json:"DateCreated,omitempty"`
	DateFirstUsage              *string    `json:"DateFirstUsage,omitempty"`
	DateSuperseded              *string    `json:"DateSuperseded,omitempty"`
	DefaultUserDeviceType       *string    `json:"DefaultUserDeviceType,omitempty"`
	DeployToMonitor             *string    `json:"DeployToMonitor,omitempty"`
	DeviceListAllowed           *string    `json:"DeviceListAllowed,omitempty"`
	ExpectedAOWattage           *int       `json:"ExpectedAOWattage,omitempty"`
	IntegratedDeviceType        *string    `json:"IntegratedDeviceType,omitempty"`
	IntegrationType             *string    `json:"IntegrationType,omitempty"`
	MergeId                     *string    `json:"MergeId,omitempty"`
	MergedDevices               *string    `json:"MergedDevices,omitempty"`
	ModelCreatedVersion         *string    `json:"ModelCreatedVersion,omitempty"`
	ModelUpdatedVersion         *string    `json:"ModelUpdatedVersion,omitempty"`
	NameUserEdit                *string    `json:"name_useredit,omitempty"`
	NameUserGuess               *string    `json:"NameUserGuess,omitempty"`
	OriginalName                *string    `json:"OriginalName,omitempty"`
	OtherSuperseded             *string    `json:"OtherSuperseded,omitempty"`
	OtherSupersededType         *string    `json:"OtherSupersededType,omitempty"`
	PeerNames                   []PeerName `json:"PeerNames,omitempty"`
	Pending                     *string    `json:"Pending,omitempty"`
	PreselectionIndex           *int       `json:"PreselectionIndex,omitempty"`
	Revoked                     *string    `json:"Revoked,omitempty"`
	SSIEnabled                  *string    `json:"SSIEnabled,omitempty"`
	SSIModel                    *string    `json:"SSIModel,omitempty"`
	SmartPlugModel              *string    `json:"SmartPlugModel,omitempty"`
	SmartPlugName               *string    `json:"SmartPlugName,omitempty"`
	TimelineAllowed             *string    `json:"TimelineAllowed,omitempty"`
	TimelineDefault             *string    `json:"TimelineDefault,omitempty"`
	Type                        *string    `json:"Type,omitempty"`
	UserControlLock             *string    `json:"UserControlLock,omitempty"`
	UserDeletable               *string    `json:"UserDeletable,omitempty"`
	UserDeleted                 *string    `json:"UserDeleted,omitempty"`
	UserDeviceType              *string    `json:"UserDeviceType,omitempty"`
	UserDeviceTypeDisplayString *string    `json:"UserDeviceTypeDisplayString,omitempty"`
	UserEditable                *string    `json:"UserEditable,omitempty"`
	UserEditableMeta            *string    `json:"UserEditableMeta,omitempty"`
	UserMergeable               *string    `json:"UserMergeable,omitempty"`
	UserShowBubble              *string    `json:"UserShowBubble,omitempty"`
	UserShowInDeviceList        *string    `json:"UserShowInDeviceList,omitempty"`
	UserSupersededBy            *string    `json:"UserSupersededBy,omitempty"`
	UserVisibleDeviceId         *string    `json:"UserVisibleDeviceId,omitempty"`
	Virtual                     *string    `json:"Virtual,omitempty"`
}

type PeerName struct {
	Name                        string  `json:"Name"`
	UserDeviceType              string  `json:"UserDeviceType"`
	Percent                     float64 `json:"Percent"`
	Icon                        string  `json:"Icon"`
	UserDeviceTypeDisplayString string  `json:"UserDeviceTypeDisplayString"`
	Make                        *string `json:"Make,omitempty"`
}

// GetDeviceByID returns the device with the given ID, or nil if not found.
func (d *Devices) GetDeviceByID(id string) *Device {
	for _, device := range d.Devices {
		if device.ID == id {
			return &device
		}
	}
	return nil
}
