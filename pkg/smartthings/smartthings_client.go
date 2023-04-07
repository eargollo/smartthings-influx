package smartthings

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const smartthingsAPI = "https://api.smartthings.com/v1"

var cli *STClient

type STClient struct {
	token string
}

func New(token string) *STClient {
	cli = &STClient{token: token}

	return cli
}

func (c STClient) Devices() (DevicesList, error) {
	var devices DevicesList

	data, err := c.get("/devices")
	if err != nil {
		return devices, err
	}

	err = json.Unmarshal(data, &devices)

	return devices, err
}

func (c STClient) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]CapabilityStatus, err error) {
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

func (c STClient) get(endpoint string) ([]byte, error) {
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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (c STClient) DeviceStatus(deviceID uuid.UUID) (status DeviceStatus, err error) {
	url := "/devices/" + deviceID.String() + "/status"

	data, err := c.get(url)
	if err != nil {
		return
	}

	// log.Printf(string(data))
	status = DeviceStatus{}

	err = json.Unmarshal([]byte(data), &status)
	return status, err
}
