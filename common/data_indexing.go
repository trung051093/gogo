package common

import (
	"encoding/json"
	"log"
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
	Data     interface{}
	Id       string
	SendTime time.Time
}

// this function to help normalize message to data index
func MessageToDataIndex(msg []byte) (*DataIndex, []byte, error) {
	var dataIndex *DataIndex
	if err := json.Unmarshal(msg, &dataIndex); err != nil {
		return nil, nil, err
	}
	dataByte, dataErr := json.Marshal(dataIndex.Data)
	if dataErr != nil {
		log.Println("Error message: ", dataErr)
	}
	return dataIndex, dataByte, nil
}
