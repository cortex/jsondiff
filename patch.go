package jsondiff

type Patch []PatchOp

func (p Patch) Apply([]byte) ([]byte, error) {
	return nil, nil
}
