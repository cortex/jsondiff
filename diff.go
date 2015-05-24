package jsondiff

import "encoding/json"
import "github.com/deckarep/golang-set"
import "errors"
import "reflect"
import "fmt"
import "path"

type PatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func DiffBytes(inb []byte, outb []byte) (ops []PatchOp, err error) {

	var in json.RawMessage
	err = json.Unmarshal(inb, &in)
	if err != nil {
		return nil, err
	}

	var out json.RawMessage
	err = json.Unmarshal(outb, &out)
	if err != nil {
		return nil, err
	}
	return DiffRaw(in, out)

}

func DiffRaw(doc1 json.RawMessage, doc2 json.RawMessage) (ops []PatchOp, err error) {
	var i1 interface{}
	var i2 interface{}
	err = json.Unmarshal(doc1, &i1)
	if err != nil {
		return
	}
	err = json.Unmarshal(doc2, &i2)
	if err != nil {
		return
	}
	return Diff(i1, i2, "/")
}

func Diff(doc1 interface{}, doc2 interface{}, root string) (ops []PatchOp, err error) {
	// If the docs have different types, just replace
	if reflect.TypeOf(doc1) != reflect.TypeOf(doc2) {
		ops = append(ops, PatchOp{"replace", "/", doc1})
		return ops, nil
	}

	// Type-specific handling
	switch t1 := doc1.(type) {
	case int:
		if t1 != doc2 {
			ops = append(ops, PatchOp{"replace", root, doc2})
		}
	case string:
		if t1 != doc2 {
			ops = append(ops, PatchOp{"replace", root, doc2})
		}
	case bool:
		if t1 != doc2 {
			ops = append(ops, PatchOp{"replace", root, doc2})
		}
	case map[string]interface{}:
		t2 := doc2.(map[string]interface{})
		extra, missing := KeyDiff(&t1, &t2)
		if extra.Cardinality() == 0 && missing.Cardinality() == 0 {
			return nil, nil
		}

		if extra.Cardinality() > 0 && missing.Cardinality() == 0 {
			for _, e := range extra.ToSlice() {
				ops = append(ops, PatchOp{"remove", path.Join(root, e.(string)), nil})
			}
			return ops, nil
		}

		if extra.Cardinality() == 0 && missing.Cardinality() > 0 {
			for _, e := range missing.ToSlice() {
				value := (t2)[e.(string)]
				ops = append(ops, PatchOp{"add", path.Join(root + e.(string)), value})
			}
			return ops, nil
		}

	case []interface{}:
		t2 := doc2.([]interface{})
		if len(t1) == len(t2) {
			for i := range t1 {
				Diff(t1[i], t2[i], path.Join(root, string(i)))
			}
		}

	default:
		return nil, errors.New(fmt.Sprintf("Not implemented: %v", reflect.TypeOf(t1)))
	}
	return
}

// KeyDiff compares the keys in two string-maps
// extra is the keys in b not in a
// missing is the keys in a not in b

func KeyDiff(a, b *map[string]interface{}) (extra, missing mapset.Set) {

	aKeys := mapset.NewSet()
	bKeys := mapset.NewSet()

	for k := range *a {
		aKeys.Add(k)
	}

	for k := range *b {
		bKeys.Add(k)
	}
	return aKeys.Difference(bKeys), bKeys.Difference(aKeys)

}
