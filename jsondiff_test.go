package jsondiff

import (
	"encoding/json"
	"fmt"
	"github.com/evanphx/json-patch"
	"os"
	"reflect"
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

type TestCase struct {
	Doc           json.RawMessage `json:"doc"`      // The JSON document to test against
	Patch         []PatchOp       `json:"patch"`    // The patch(es) to apply
	Expected      json.RawMessage `json:"expected"` // The expected resulting document, OR
	ExpectedError string          `json:"error"`    // A string describing an expected error
	Comment       string          `json:"comment"`  // A string describing the test
	Disabled      bool            `json:"disabled"` // True if the test should be skipped
}

func TestSuite(t *testing.T) {
	var ok, failed, errors, skipped int

	testf, err := os.Open("tests/spec_tests.json")
	if err != nil {
		t.Fatalf("Failed to open test suite")
	}
	dec := json.NewDecoder(testf)
	var testCases []TestCase
	err = dec.Decode(&testCases)
	if err != nil {
		t.Fatalf("Failed to parse test cases: %v", err)
	}
	for _, test := range testCases {
		if test.ExpectedError != "" {
			t.Log("== SKIP: ", test.Comment)
			skipped++
			continue
		}
		t.Log("== RUN: ", test.Comment)
		t.Log("In: ", string(test.Doc))
		t.Log("Out: ", string(test.Expected))
		t.Log("Expected patch: ", test.Patch)

		patch, err := DiffRaw(test.Doc, test.Expected)
		if err != nil {
			t.Log("Error: ", err)
			t.Error(test.Comment, err)
			t.Log("-- ERROR\n")
			errors++
			continue
		}
		t.Log("Generated patch", test.Patch)

		if !(reflect.DeepEqual(patch, test.Patch) ||
			len(patch) == len(test.Patch) && len(patch) == 0) {
			t.Logf("Failed: not equal: %v %v", patch, test.Patch)
			t.Fail()
			t.Log("-- FAIL\n")
			failed++
			continue
		}
		fmt.Println(test.Comment, patch)
		t.Log("-- OK\n")
		ok++
	}
	t.Logf("OK: %v Failed: %v Error: %v Skipped: %v", ok, failed, errors, skipped)
}

func verifyPatch(t *testing.T, in []byte, out []byte, expected []byte) {

	patchObj, err := DiffBytes(in, out)
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
