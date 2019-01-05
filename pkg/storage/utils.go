package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Exists checks if a file/folder exists
func Exists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// SaveJSON saves the object to a JSON file
func SaveJSON(filename string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

// ReadJSON reads the object from a JSON file
func ReadJSON(filename string, v interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
