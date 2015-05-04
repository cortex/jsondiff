package jsondiff

import "encoding/json"
import "errors"
import "github.com/evanphx/json-patch"

func Diff(inb []byte, outb []byte) ([]byte, error) {
	if jsonpatch.Equal(inb, outb) {
		return []byte("[]"), nil
	}
	var in interface{}
	err := json.Unmarshal(inb, &in)
	if err != nil {
		return nil, err
	}

	var out interface{}
	err = json.Unmarshal(outb, &out)
	if err != nil {
		return nil, err
	}

	switch in.(type) {
	case map[string]interface{}:
		return nil, errors.New("Not implemented: ")
	default:
		return nil, errors.New("Not Implemented")
	}
}
