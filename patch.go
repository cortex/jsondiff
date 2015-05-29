package jsondiff

import "encoding/json"

type Patch []PatchOp

func NewPatchFromBytes(b []byte) (*Patch, error) {
	p := Patch{}
	err := json.Unmarshal(b, &b)
	return &p, err
}

func (p Patch) Apply([]byte) ([]byte, error) {
	return nil, nil
}
