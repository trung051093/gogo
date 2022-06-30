package common

import (
	"encoding/json"
	"errors"
	"log"
)

func JsonSprint[T any](msg string, data T) {
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	log.Println(msg + string(out))
}

func CompactJson[T any](data T) T {
	var dataJs T = data

	dataBytes, _ := json.Marshal(data)

	if json.Valid(dataBytes) {
		json.Unmarshal(dataBytes, &dataJs)
	}

	return dataJs
}

func JsonToString[T any](data T) (string, error) {
	dataBytes, _ := json.Marshal(data)
	if json.Valid(dataBytes) {
		str := string(dataBytes)
		return str, nil
	}
	return "", errors.New("invalid json string")
}

func StringToJson[T any](str string, data T) error {
	if err := json.Unmarshal([]byte(str), data); err != nil {
		return err
	}
	return nil
}

func JsonToByte[T any](data T) []byte {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return dataByte
}
