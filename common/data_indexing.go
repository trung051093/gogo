package common

import (
	"bytes"
	"encoding/json"
	"time"
)

type Action string

const (
	Update Action = "update"
	Delete Action = "delete"
	Create Action = "create"
)

type DataIndex struct {
	Index    string
	Action   Action
	Data     any
	Id       string
	SendTime time.Time
}

// normalize message to data index
func (d *DataIndex) GetByte() ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBytes).Encode(d); err != nil {
		return nil, err
	}
	return reqBodyBytes.Bytes(), nil
}

// normalize message to data index
func (d *DataIndex) FromByte(msg []byte) error {
	if err := json.Unmarshal(msg, d); err != nil {
		return err
	}
	return nil
}
