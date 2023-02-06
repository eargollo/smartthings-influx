package smartthings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const smartthingsAPI = "https://api.smartthings.com/v1"

type SmartThingsTransport struct {
	token string
}

func NewTransport(token string) *SmartThingsTransport {
	return &SmartThingsTransport{token: token}
}

func (c SmartThingsTransport) Devices() (DevicesList, error) {
	var devices DevicesList

	data, err := c.get("/devices")
	if err != nil {
		return devices, err
	}

	err = json.Unmarshal(data, &devices)

	for i, dev := range devices.Items {
		dev.Client = &c
		devices.Items[i] = dev
	}

	return devices, err
}

func (c SmartThingsTransport) DeviceStatus(deviceID uuid.UUID) (DeviceStatus, error) {
	var status DeviceStatus

	url := "/devices/" + deviceID.String() + "/status"

	data, err := c.get(url)
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(data, &status)

	return status, err
}

func (c SmartThingsTransport) get(endpoint string) ([]byte, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", smartthingsAPI+endpoint, nil)
	if err != nil {
		return []byte{}, err
	}

	// add authorization header to the req
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (c SmartThingsTransport) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]CapabilityStatus, err error) {
	url := "/devices/" + deviceID.String() + "/components/" + componentId + "/capabilities/" + capabilityId + "/status"

	data, err := c.get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &status)
	if err != nil {
		return status, fmt.Errorf("could not unmarshall device capability status payload: '%s'", string(data))
	}

	return status, err
}
