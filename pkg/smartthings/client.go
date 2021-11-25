package smartthings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const api = "https://api.smartthings.com/v1"

var cli *Client

type Client struct {
	token string
}

func Init(token string) *Client {
	cli = &Client{token: token}
	return cli
}

func Devices() (DevicesList, error) {
	return cli.Devices()
}

func (c Client) Devices() (devices DevicesList, err error) {
	data, err := c.get("/devices")
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(data), &devices)

	return
}

func (c Client) DeviceStatus(deviceID uuid.UUID) (status DeviceStatus, err error) {
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

func (c Client) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]interface{}, err error) {
	url := "/devices/" + deviceID.String() + "/components/" + componentId + "/capabilities/" + capabilityId + "/status"

	log.Printf("Endpoing is '%s'", url)
	data, err := c.get(url)
	if err != nil {
		return
	}

	log.Printf(string(data))
	err = json.Unmarshal([]byte(data), &status)
	return status, err
}

func (c Client) get(endpoint string) ([]byte, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", api+endpoint, nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
