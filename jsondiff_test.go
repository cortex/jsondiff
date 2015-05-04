package jsondiff

import (
	"fmt"
	"github.com/evanphx/json-patch"
	"testing"
)

func TestEqual(t *testing.T) {
	a := []byte(`{}`)
	b := []byte(`{}`)
	expected := []byte(`[]`)
	verify(t, a, b, expected)
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
	verify(t, a, b, expected)
}

func verify(t *testing.T, in []byte, out []byte, expected []byte) {

	patch, err := Diff(in, out)
	if err != nil {
		t.Error(err)
	}
	obj, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		t.Error(err)
	}
	result, err := obj.Apply(in)
	if err != nil {
		t.Error(err)
	}

	if !jsonpatch.Equal(out, result) {
		fmt.Printf("in: %s out: %s", out, result)
		t.Fail()
	}

	if !jsonpatch.Equal(expected, patch) {
		fmt.Printf("expected: %s out: %s", expected, patch)
		t.Fail()
	}
}
