package main

import (
	"bytes"
	"encoding/gob"

	apiv1 "github.com/loshz/platform/internal/api/v1"
)

// TODO: add more event types.
func GenerateRandomEvent(machineId string) (*apiv1.SendEventRequest, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&apiv1.NetworkEvent{}); err != nil {
		return nil, err
	}

	req := &apiv1.SendEventRequest{
		MachineId: machineId,
		Type:      apiv1.EventType_EVENT_TYPE_NETWORK,
		Data:      buf.Bytes(),
	}

	return req, nil
}
