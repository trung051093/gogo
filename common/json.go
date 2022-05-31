package common

import (
	"encoding/json"
	"errors"
	"log"
)

func JsonSprint(msg string, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	log.Println(msg + string(out))
}

func CompactJson(data interface{}) interface{} {
	var dataJs interface{} = data

	dataBytes, _ := json.Marshal(data)

	if json.Valid(dataBytes) {
		json.Unmarshal(dataBytes, &dataJs)
	}

	return dataJs
}

func JsonToString(data interface{}) (string, error) {
	dataBytes, _ := json.Marshal(data)
	if json.Valid(dataBytes) {
		str := string(dataBytes)
		return str, nil
	}
	return "", errors.New("Invalid json string")
}

func StringToJson(str string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func JsonToByte(data interface{}) []byte {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil
	}
	return dataByte
}
