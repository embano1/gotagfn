package faastagger

import (
	"time"

	"github.com/vmware/govmomi/vim25/types"
)

// InbountEvent is the JSON object forwarded by OpenFaaS to this function
type InbountEvent struct {
	Topic    string `json:"topic,omitempty"`
	Category string `json:"category,omitempty"`

	UserName    string                        `json:"userName,omitempty"`
	CreatedTime time.Time                     `json:"createdTime,omitempty"`
	Object      string                        `json:"object,omitempty"`
	UUID        string                        `json:"uuid,omitempty"`
	MoRef       *types.ManagedObjectReference `json:"moref,omitempty"`
}
