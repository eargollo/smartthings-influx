package smartthings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const api = "https://api.smartthings.com/v1"

var cli *Client

type Client struct {
	token         string
	conversionMap map[string]map[string]float64
}

func Init(token string, conversionMap map[string]map[string]float64) *Client {
	cli = &Client{token: token, conversionMap: conversionMap}

	return cli
}

func ParseConversionMap(valuemap map[string]interface{}) (map[string]map[string]float64, error) {
	conversionMap := make(map[string]map[string]float64)
	for key, value := range valuemap {
		innermap, ok := value.(map[string]any)
		if !ok {
			return map[string]map[string]float64{}, fmt.Errorf("could not parse valuemap, it shold be in the format of metric to a map of values, got %v", value)
		}

		_, ok = conversionMap[key]
		if !ok {
			conversionMap[key] = make(map[string]float64)
		}

		for inkey, inval := range innermap {
			var inFloat float64
			_, ok := inval.(int)
			if ok {
				inFloat = float64(inval.(int))
			} else {
				inFloat, ok = inval.(float64)
				if !ok {
					return map[string]map[string]float64{}, fmt.Errorf("could not convert %v to a number for metric %s", inval, key)
				}
			}

			list := conversionMap[key]
			list[inkey] = inFloat
			conversionMap[key] = list
		}
	}
	return conversionMap, nil
}

func Devices() (DevicesList, error) {
	return cli.Devices()
}

func (c Client) ConvertValueToFloat(metric string, value any) (float64, error) {
	_, ok := value.(float64)
	if ok {
		return value.(float64), nil
	}

	_, ok = value.(string)
	if ok {
		stValue := value.(string)
		// Check if there is a map for metric
		metricMap, ok := c.conversionMap[metric]
		if !ok {
			return 0, fmt.Errorf("there is no value map for metric '%s' and value '%s', can't convert", metric, stValue)
		}
		return metricMap[stValue], nil
	}
	return 0, nil
}

func (c Client) Devices() (devices DevicesList, err error) {
	data, err := c.get("/devices")
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(data), &devices)
	devices.client = &c

	return
}

func (c Client) DevicesWithCapabilities(capabilities []string) (list DevicesWithCapabilitiesResult, err error) {
	data, err := c.Devices()
	if err != nil {
		return
	}

	list = data.DevicesWithCapabilities(capabilities)
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

type CapabilityStatus struct {
	Timestamp time.Time `json:"timestamp"`
	Unit      string    `json:"unit"`
	Value     any       `json:"value"`
}

func (status CapabilityStatus) FloatValue(metric string) (float64, error) {
	return cli.ConvertValueToFloat(metric, status.Value)
}

func (c Client) DeviceCapabilityStatus(deviceID uuid.UUID, componentId string, capabilityId string) (status map[string]CapabilityStatus, err error) {
	url := "/devices/" + deviceID.String() + "/components/" + componentId + "/capabilities/" + capabilityId + "/status"

	// log.Printf("Endpoint is '%s'", url)
	data, err := c.get(url)
	if err != nil {
		return
	}

	// log.Printf("Status for device '%s' component '%s' capability '%s' payload is '%s' ",
	// 	deviceID,
	// 	componentId,
	// 	capabilityId,
	// 	string(data),
	// )

	err = json.Unmarshal(data, &status)
	if err != nil {
		return status, fmt.Errorf("could not unmarshall device capability status payload: '%s'", string(data))
	}

	// log.Printf("Unmarshalled status for device '%s' component '%s' capability '%s' payload is '%v' ",
	// 	deviceID,
	// 	componentId,
	// 	capabilityId,
	// 	status,
	// )

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
