package jsondiff

import "encoding/json"
import "github.com/deckarep/golang-set"
import "errors"

type Dict map[string]interface{}

type PatchOp struct {
	Op    string  `json:"op"`
	Path  string  `json:"path"`
	Value *string `json:"value,omitempty"`
}

func Diff(inb []byte, outb []byte) (ops []PatchOp, err error) {

	var in Dict
	err = json.Unmarshal(inb, &in)
	if err != nil {
		return nil, err
	}

	var out Dict
	err = json.Unmarshal(outb, &out)
	if err != nil {
		return nil, err
	}

	extra, missing := KeyDiff(in, out)

	if extra.Cardinality() == 0 && missing.Cardinality() == 0 {
		return nil, nil
	}

	if extra.Cardinality() > 0 && missing.Cardinality() == 0 {
		for _, e := range extra.ToSlice() {
			ops = append(ops, PatchOp{"remove", "/" + e.(string), nil})
		}
		return ops, nil
	}

	if extra.Cardinality() == 0 && missing.Cardinality() > 0 {
		for _, e := range missing.ToSlice() {
			value := out[e.(string)]
			v, _ := json.Marshal(value)
			vs := string(v[1 : len(v)-1]) // FIXME: only works for strings
			ops = append(ops, PatchOp{"add", "/" + e.(string), &vs})
		}
		return ops, nil
	}

	// TODO:
	// if len extra > 0 and len (missing) > 0
	// add and remove
	// future: rename?

	return nil, errors.New("Not implemented")
}

// KeyDiff compares the keys in two string-maps
// extra is the keys in b not in a
// missing is the keys in a not in b

func KeyDiff(a, b Dict) (extra, missing mapset.Set) {

	aKeys := mapset.NewSet()
	bKeys := mapset.NewSet()

	for k := range a {
		aKeys.Add(k)
	}

	for k := range b {
		bKeys.Add(k)
	}
	return aKeys.Difference(bKeys), bKeys.Difference(aKeys)

}
