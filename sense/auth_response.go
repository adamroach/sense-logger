package sense

import "time"

type AuthResponse struct {
	Authorized   bool          `json:"authorized"`
	AccountID    int           `json:"account_id"`
	UserID       int           `json:"user_id"`
	AccessToken  string        `json:"access_token"`
	Roles        string        `json:"roles"` // only in the response from the refresh token endpoint
	RefreshToken string        `json:"refresh_token"`
	BridgeServer string        `json:"bridge_server"` // only in the response from the login endpoint
	DateCreated  string        `json:"date_created"`  // only in the response from the login endpoint
	TotpEnabled  bool          `json:"totp_enabled"`
	Settings     UserSettings  `json:"settings"` // only in the response from the login endpoint
	Monitors     []MonitorInfo `json:"monitors"` // only in the response from the login endpoint
	Expires      time.Time     `json:"expires"`  // only in the response from the refresh token endpoint
}

type UserSettings struct {
	UserID   int             `json:"user_id"`
	Settings SettingsDetails `json:"settings"`
	Version  int             `json:"version"`
}

type SettingsDetails struct {
	Timeline             TimelineSettings     `json:"timeline"`
	Alerts               AlertSettings        `json:"alerts"`
	Notifications        NotificationSettings `json:"notifications"`
	LabsEnabled          bool                 `json:"labs_enabled"`
	HideTrendsCarbonCard bool                 `json:"hide_trends_carbon_card"`
}

type TimelineSettings struct {
	ShowDevices map[string]map[string]bool `json:"show_devices"`
}

type AlertSettings struct {
	Enabled map[string]map[string]bool `json:"enabled"`
}

type NotificationSettings map[string]map[string]bool

type MonitorInfo struct {
	ID                         int               `json:"id"`
	DateCreated                string            `json:"date_created"`
	SerialNumber               string            `json:"serial_number"`
	TimeZone                   string            `json:"time_zone"`
	SolarConnected             bool              `json:"solar_connected"`
	SolarConfigured            bool              `json:"solar_configured"`
	Online                     bool              `json:"online"`
	Attributes                 MonitorAttributes `json:"attributes"`
	SignalCheckCompletedTime   string            `json:"signal_check_completed_time"`
	DataSharing                []string          `json:"data_sharing"`
	EthernetSupported          bool              `json:"ethernet_supported"`
	PowerOverEthernetSupported bool              `json:"power_over_ethernet_supported"`
	AuxIgnore                  bool              `json:"aux_ignore"`
	AuxPort                    string            `json:"aux_port"`
	HardwareType               string            `json:"hardware_type"`
	ZigbeeSupported            bool              `json:"zigbee_supported"`
}

type MonitorAttributes struct {
	ID                     int     `json:"id"`
	Name                   string  `json:"name"`
	State                  string  `json:"state"`
	Cost                   float64 `json:"cost"`
	SellBackRate           float64 `json:"sell_back_rate"`
	UserSetCost            bool    `json:"user_set_cost"`
	CycleStart             int     `json:"cycle_start"`
	BasementType           string  `json:"basement_type"`
	HomeSizeType           string  `json:"home_size_type"`
	HomeType               string  `json:"home_type"`
	NumberOfOccupants      string  `json:"number_of_occupants"`
	OccupancyType          string  `json:"occupancy_type"`
	YearBuiltType          string  `json:"year_built_type"`
	BasementTypeKey        *string `json:"basement_type_key"`
	HomeSizeTypeKey        *string `json:"home_size_type_key"`
	HomeTypeKey            *string `json:"home_type_key"`
	OccupancyTypeKey       *string `json:"occupancy_type_key"`
	YearBuiltTypeKey       *string `json:"year_built_type_key"`
	Address                *string `json:"address"`
	City                   *string `json:"city"`
	PostalCode             string  `json:"postal_code"`
	ElectricityCost        *string `json:"electricity_cost"`
	ShowCost               bool    `json:"show_cost"`
	TOUEnabled             bool    `json:"tou_enabled"`
	SolarTOUEnabled        bool    `json:"solar_tou_enabled"`
	PowerRegion            *string `json:"power_region"`
	ToGridThreshold        *string `json:"to_grid_threshold"`
	Panel                  *string `json:"panel"`
	HomeInfoSurveyProgress string  `json:"home_info_survey_progress"`
	DeviceSurveyProgress   string  `json:"device_survey_progress"`
	UserSetSellBackRate    bool    `json:"user_set_sell_back_rate"`
}
