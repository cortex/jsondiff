package jsondiff

import (
	"encoding/json"
	"fmt"
	"github.com/evanphx/json-patch"
	"testing"
)

func TestEqual(t *testing.T) {
	a := []byte(`{}`)
	b := []byte(`{}`)
	expected := []byte(`[]`)
	verifyPatch(t, a, b, expected)
}

func TestSimpleAdd(t *testing.T) {
	a := []byte(`{ "foo": "bar"}`)

	b := []byte(`{
     "baz": "qux",
     "foo": "bar"
   }`)

	expected := []byte(`[
     { "op": "add", "path": "/baz", "value": "qux" }
   ]`)
	verifyPatch(t, a, b, expected)
}

func TestSimpleRemove(t *testing.T) {
	a := []byte(`{
     "baz": "qux",
     "foo": "bar"
   }`)

	b := []byte(`{ "foo": "bar"}`)

	expected := []byte(`[
     { "op": "remove", "path": "/baz"}
   ]`)
	verifyPatch(t, a, b, expected)
}
func verifyPatch(t *testing.T, in []byte, out []byte, expected []byte) {

	patchObj, err := Diff(in, out)
	if err != nil {
		t.Error(err)
	}
	patch, _ := json.Marshal(patchObj)
	obj, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		t.Error(err)
	}
	result, err := obj.Apply(in)
	if err != nil {
		t.Error(err)
	}

	if !jsonpatch.Equal(out, result) {
		fmt.Printf("in: %s\n", out)
		fmt.Printf("out: %s\n", result)
		t.Fail()
	}

	if !jsonpatch.Equal(expected, patch) {
		fmt.Printf("expected patch:\n %s \nactual patch:\n %s\n", expected, patch)
		t.Fail()
	}
}
