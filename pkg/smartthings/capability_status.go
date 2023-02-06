package smartthings

import (
	"time"
)

type CapabilityStatus struct {
	Timestamp time.Time `json:"timestamp"`
	Unit      string    `json:"unit"`
	Value     any       `json:"value"`
}
