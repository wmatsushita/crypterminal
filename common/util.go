package common

import (
	"encoding/json"
	"io/ioutil"
)

func LoadFromJsonFile(fileName string, v interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return NewErrorWithCause("Could not read Json config file.", err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return NewErrorWithCause("Could not deserialize json file.", err)
	}

	return nil
}

func Contains(values []string, value string) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}
	return false
}
