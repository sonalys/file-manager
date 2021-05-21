package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

func read(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("couldn't read the file %s", path))
	}
	return content, nil
}

// ReadJSON reads an JSON content into a structure, this structure should be an pointer.
func ReadJSON(path string, dest interface{}) error {
	content, err := read(path)
	if err != nil {
		return errors.Wrap(err, "couldn't read the JSON file")
	}

	if err := json.Unmarshal(content, dest); err != nil {
		return errors.Wrap(err, "couldn't unmarshal the JSON file")
	}
	return nil
}
