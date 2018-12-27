package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Exists checks if a file/folder exists
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// SaveJSON creates/overwrites a JSON file with the specified object data
func SaveJSON(fileName string, object interface{}) error {
	data, err := json.MarshalIndent(object, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, 0644)
}

// ReadJSON get the object data from a JSON file
func ReadJSON(fileName string, object interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}
