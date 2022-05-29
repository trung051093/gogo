package common

import (
	"encoding/json"
	"errors"
)

func CompactJson(data interface{}) interface{} {
	var dataJs interface{} = data

	dataBytes, _ := json.Marshal(data)

	if json.Valid(dataBytes) {
		json.Unmarshal(dataBytes, &dataJs)
	}

	return dataJs
}

func ConvertJsonToString(data interface{}) (string, error) {
	dataBytes, _ := json.Marshal(data)
	if json.Valid(dataBytes) {
		str := string(dataBytes)
		return str, nil
	}
	return "", errors.New("Invalid json string")
}

func ConvertStringToJson(str string) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, err
	}
	return data, nil
}
