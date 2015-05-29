package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)
import "os"
import "log"
import "text/template"
import "regexp"
import "flag"
import "fmt"

type TestCase struct {
	Doc           json.RawMessage `json:"doc"`      // The JSON document to test against
	Patch         json.RawMessage `json:"patch"`    // The patch(es) to apply
	Expected      json.RawMessage `json:"expected"` // The expected resulting document, OR
	ExpectedError string          `json:"error"`    // A string describing an expected error
	Comment       string          `json:"comment"`  // A string describing the test
	Disabled      bool            `json:"disabled"` // True if the test should be skipped
	Name          string
}

const headerTemplate = `package {{.Package}}

import "testing"

`

const testTemplate = `
func Test{{.Name}} (t* testing.T){
	//{{.Comment}}
	doc := []byte(` + "`" + `{{printf "%s" .Doc }}` + "`)" + `
	patch := []byte(` + "`" + `{{printf "%s" .Patch}}` + "`)" + `
{{if not .ExpectedError}}	expected := ` + "[]byte(`" + `{{printf "%s" .Expected}}` + "`)" + `
	verifyPatch(t, doc, patch, expected){{else}}expectedError := "{{.ExpectedError}}"
	verifyPatchError(t, doc, patch, expectedError)
{{end}}
}
`

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func CamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return string(bytes.Join(chunks, nil))
}

func testName(in string) string {
	re := regexp.MustCompile("[^0-9A-Za-z ]")
	return strings.Title(CamelCase(re.ReplaceAllString(in, "")))
}

func main() {
	outfile := flag.String("o", "", "output filename")
	flag.Parse()

	if *outfile == "" {
		fmt.Println("Outfile is required")
		os.Exit(1)
	}

	name := flag.Arg(0)
	fmt.Println("outfile:", *outfile)

	outf, err := os.Create(*outfile)
	if err != nil {
		panic(err)
	}
	defer outf.Close()
	tmpl, err := template.New("TestCase").Parse(testTemplate)
	if err != nil {
		panic(err)
	}

	htmpl, err := template.New("TestCase").Parse(headerTemplate)
	if err != nil {
		panic(err)
	}

	testf, err := os.Open(name)
	if err != nil {
		log.Fatalf("Failed to open test suite")
	}

	dec := json.NewDecoder(testf)
	var testCases []TestCase
	err = dec.Decode(&testCases)

	if err != nil {
		log.Fatalf("Failed to parse test cases: %v", err)
	}
	htmpl.Execute(outf, struct{ Package string }{os.Getenv("GOPACKAGE")})
	for i, test := range testCases {
		test.Name = testName(test.Comment)
		if test.Name == "" {
			test.Name = strconv.Itoa(i)
		}
		err = tmpl.Execute(outf, test)
	}
}
