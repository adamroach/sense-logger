package sense

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// DefaultHost is the default host for the Sense API
	ApiHost                = "api.sense.com"
	WebsocketHost          = "clientrt.sense.com"
	AuthEndpoint           = "https://" + ApiHost + "/apiservice/api/v1/authenticate"
	RefreshEndpoint        = "https://" + ApiHost + "/apiservice/api/v1/renew"
	LabsTemplate           = "https://{{.ApiHost}}/apiservice/api/v1/app/monitors/{{.MonitorID}}/labs_content"
	DeviceOverviewTemplate = "https://{{.ApiHost}}/apiservice/api/v1/app/monitors/{{.MonitorID}}/devices/overview?" +
		"include_merged={{.IncludeMerged}}"
	DeviceDetailTemplate = "https://{{.ApiHost}}/apiservice/api/v1/app/monitors/{{.MonitorID}}/devices/{{.DeviceID}}"
	WebsocketTemplate    = "wss://{{.WebsocketHost}}/monitors/{{.MonitorID}}/realtimefeed?" +
		"access_token={{.AccessToken}}&" +
		"sense_device_id={{.DeviceID}}&" +
		"sense_protocol_version={{.ProtocolVersion}}&" +
		"sense_client_type={{.ClientType}}&" +
		"sense_ui_language={{.UILanguage}}"
	defaultTokenExpiry = 24 * time.Hour
	watchdogInterval   = 5 * time.Second
)

type WebsocketParams struct {
	WebsocketHost   string
	MonitorID       int
	AccessToken     string
	DeviceID        string
	ProtocolVersion int
	ClientType      string
	UILanguage      string
}

type DeviceOverviewParams struct {
	ApiHost       string
	MonitorID     int
	IncludeMerged bool
}

type DeviceDetailParams struct {
	ApiHost   string
	MonitorID int
	DeviceID  string
}

type LabsReportParams struct {
	ApiHost   string
	MonitorID int
}

var (
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")
	ErrNotAuthenticated     = fmt.Errorf("not authenticated")
	ErrNoDevicesLoaded      = fmt.Errorf("no devices loaded")
	ErrDeviceNotFound       = fmt.Errorf("device not found")
	ErrFailedToLoadDevices  = fmt.Errorf("failed to load devices")
)

type Client struct {
	client       *http.Client
	clientId     string
	authResponse *AuthResponse
	devices      *Devices
	updates      chan *RealtimeUpdate
	watchdog     *time.Timer
	conn         *websocket.Conn
	mu           sync.Mutex
}

func NewClient() *Client {
	return &Client{
		client:   &http.Client{Timeout: 10 * time.Second},
		clientId: generateRandomClientID(128),
	}
}

func (c *Client) Login(username, password string) error {
	form := "email=" + username + "&password=" + password
	req, err := http.NewRequest("POST", AuthEndpoint, strings.NewReader(form))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.authResponse = &AuthResponse{
		Expires: time.Now().Add(defaultTokenExpiry),
	}
	if err := json.NewDecoder(resp.Body).Decode(&c.authResponse); err != nil {
		return err
	}
	if !c.authResponse.Authorized {
		return ErrAuthenticationFailed
	}
	return c.loadDevices()
}

func (c *Client) Refresh() (until time.Time, err error) {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return time.Time{}, ErrNotAuthenticated
	}
	form := fmt.Sprintf(
		"user_id=%d&is_access_token=true&refresh_token=%s",
		c.authResponse.UserID,
		c.authResponse.RefreshToken,
	)
	req, err := http.NewRequest("POST", RefreshEndpoint, strings.NewReader(form))
	if err != nil {
		return time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	// We overwrite the authResponse with the new one, leaving any fields not present in the new response as-is
	c.authResponse.Expires = time.Now().Add(defaultTokenExpiry)
	if err := json.NewDecoder(resp.Body).Decode(&c.authResponse); err != nil {
		return time.Time{}, err
	}
	if !c.authResponse.Authorized {
		return time.Time{}, ErrAuthenticationFailed
	}
	return c.authResponse.Expires, nil
}

func (c *Client) TokenExpiry() time.Time {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return time.Time{}
	}
	return c.authResponse.Expires
}

func (c *Client) Close() error {
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	if c.watchdog != nil {
		c.watchdog.Stop()
		c.watchdog = nil
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) loadDevices() error {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return ErrNotAuthenticated
	}

	params := DeviceOverviewParams{
		ApiHost:       ApiHost,
		MonitorID:     c.authResponse.Monitors[0].ID,
		IncludeMerged: true,
	}

	tmpl, err := template.New("deviceOverviewEndpoint").Parse(DeviceOverviewTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse DeviceOverview template: %w", err)
	}

	var deviceOverviewEndpointBuilder strings.Builder
	if err := tmpl.Execute(&deviceOverviewEndpointBuilder, params); err != nil {
		return fmt.Errorf("failed to execute DeviceOverview template: %w", err)
	}

	DeviceOverviewEndpoint := deviceOverviewEndpointBuilder.String()

	req, err := http.NewRequest("GET", DeviceOverviewEndpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "bearer "+c.authResponse.AccessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", ErrFailedToLoadDevices, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.devices = &Devices{}
	if err := json.Unmarshal(bodyBytes, c.devices); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDevices() (*Devices, error) {
	if c.devices == nil {
		return nil, ErrNoDevicesLoaded
	}
	return c.devices, nil
}

func (c *Client) GetDeviceByID(id string) (*Device, error) {
	if c.devices == nil {
		return nil, ErrNoDevicesLoaded
	}
	device := c.devices.GetDeviceByID(id)
	if device == nil {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}

func (c *Client) GetDeviceDetails(deviceID string) (*DeviceDetails, error) {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return nil, ErrNotAuthenticated
	}

	params := DeviceDetailParams{
		ApiHost:   ApiHost,
		MonitorID: c.authResponse.Monitors[0].ID,
		DeviceID:  deviceID,
	}

	tmpl, err := template.New("deviceDetailEndpoint").Parse(DeviceDetailTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DeviceDetail template: %w", err)
	}

	var deviceDetailEndpointBuilder strings.Builder
	if err := tmpl.Execute(&deviceDetailEndpointBuilder, params); err != nil {
		return nil, fmt.Errorf("failed to execute DeviceDetail template: %w", err)
	}

	deviceDetailEndpoint := deviceDetailEndpointBuilder.String()

	req, err := http.NewRequest("GET", deviceDetailEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+c.authResponse.AccessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch device details: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	deviceDetails := &DeviceDetails{}
	if err := json.Unmarshal(bodyBytes, deviceDetails); err != nil {
		return nil, err
	}

	return deviceDetails, nil
}

func (c *Client) GetRealtimeUpdate() (*RealtimeUpdate, error) {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return nil, ErrNotAuthenticated
	}
	c.mu.Lock()
	if c.conn == nil {
		if err := c.startRealtimeUpdates(); err != nil {
			c.mu.Unlock()
			return nil, err
		}
		c.watchdog = time.AfterFunc(watchdogInterval, func() { c.Close() })
	}
	c.mu.Unlock()
	update := <-c.updates
	if update == nil {
		return nil, io.EOF
	}
	for i := range update.Payload.Devices {
		if update.Payload.Devices[i].ID == "" {
			continue
		}
		device, err := c.GetDeviceByID(update.Payload.Devices[i].ID)
		if err != nil {
			continue
		}
		// Update the device with the latest data from the realtime update
		device.Attrs = update.Payload.Devices[i].Attrs
		device.Watts = update.Payload.Devices[i].Watts
		device.Cirumference = update.Payload.Devices[i].Cirumference
		device.StatusDetails = update.Payload.Devices[i].StatusDetails
		device.AlwaysOnState = update.Payload.Devices[i].AlwaysOnState
		device.AlwaysOnWatts = update.Payload.Devices[i].AlwaysOnWatts
		update.Payload.Devices[i] = *device
	}

	return update, nil
}

func (c *Client) GetLabsReport() (*LabsReport, error) {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return nil, ErrNotAuthenticated
	}

	params := LabsReportParams{
		ApiHost:   ApiHost,
		MonitorID: c.authResponse.Monitors[0].ID,
	}

	tmpl, err := template.New("labsReportEndpoint").Parse(LabsTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LabsReport template: %w", err)
	}

	var labsReportEndpointBuilder strings.Builder
	if err := tmpl.Execute(&labsReportEndpointBuilder, params); err != nil {
		return nil, fmt.Errorf("failed to execute LabsReport template: %w", err)
	}

	labsReportEndpoint := labsReportEndpointBuilder.String()

	req, err := http.NewRequest("GET", labsReportEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+c.authResponse.AccessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch labs report: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	labsReport := &LabsReport{}
	if err := json.Unmarshal(bodyBytes, labsReport); err != nil {
		return nil, err
	}

	return labsReport, nil
}

func (c *Client) startRealtimeUpdates() error {
	if c.authResponse == nil || !c.authResponse.Authorized {
		return ErrNotAuthenticated
	}

	c.updates = make(chan *RealtimeUpdate, 1024)

	params := WebsocketParams{
		WebsocketHost:   WebsocketHost,
		MonitorID:       c.authResponse.Monitors[0].ID,
		AccessToken:     c.authResponse.AccessToken,
		DeviceID:        c.clientId,
		ProtocolVersion: 11,
		ClientType:      "web",
		UILanguage:      "en-US",
	}

	tmpl, err := template.New("websocketURI").Parse(WebsocketTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse WebSocket template: %w", err)
	}

	var websocketURIBuilder strings.Builder
	if err := tmpl.Execute(&websocketURIBuilder, params); err != nil {
		return fmt.Errorf("failed to execute WebSocket template: %w", err)
	}

	websocketURI := websocketURIBuilder.String()
	c.conn, _, err = websocket.DefaultDialer.Dial(websocketURI, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	go func() {
		defer func() {
			c.mu.Lock()
			close(c.updates)
			c.conn.Close()
			c.conn = nil
			c.mu.Unlock()
		}()
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("WebSocket read error: %v\n", err)
				}
				return
			}

			update := &RealtimeUpdate{}
			if err := json.Unmarshal(message, update); err != nil {
				log.Printf("Failed to unmarshal WebSocket message: %v\n", err)
				continue
			}

			if update.Type != "realtime_update" {
				log.Printf("Unexpected WebSocket message type: %s\n", update.Type)
				continue
			}

			if len(update.Payload.Devices) > 0 {
				c.watchdog.Reset(watchdogInterval)
			}

			c.updates <- update
		}
	}()

	return nil
}

func generateRandomClientID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			panic("failed to generate random byte")
		}
		result[i] = charset[randomByte[0]%byte(len(charset))]
	}
	return string(result)
}
