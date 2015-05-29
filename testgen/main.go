package main

import "encoding/json"
import "os"
import "log"
import "text/template"
import "strconv"
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

const testTemplate = `
//{{.Comment}}
func Test{{.Name}} (t* testing.T){
    doc := ` + "`" + `{{printf "%s" .Doc }}` + "`" + `
    patch := ` + "`" + `{{printf "%s" .Patch}}` + "`" + `
{{if .ExpectedError}}
    expected := ` + "`" + `{{printf "%s" .Expected}}` + "`" + `
    verifyPatch(doc, patch, expected)
}
`

func main() {
	outfile := flag.String("o", "", "output filename")
	flag.Parse()
	name := flag.Arg(0)
	fmt.Println(*outfile)
	outf, err := os.Create(*outfile)
	if err != nil {
		panic(err)
	}
	defer outf.Close()
	tmpl, err := template.New("TestCase").Parse(testTemplate)
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
	for i, test := range testCases {
		test.Name = strconv.Itoa(i)
		fmt.Println(test.Name)
		err = tmpl.Execute(outf, test)
	}
}
