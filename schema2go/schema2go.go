//go:generate go tool yacc items.y

package main

import (
	"flag"
	"io"
	"io/ioutil"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"log"
	"strings"
	"text/template"
)

const temp = `package {{ .Package }}
/* 
autogenerated from:
schema files:
{{ range $v := .SchemaFiles }}{{$v}} {{end}}

object description:
{{ .ObjectDesc }}
*/

import (
	"fmt"
	"errors"
	"strings"
	"github.com/bytemine/ldap-crud/crud"
	"github.com/rbns/ldap"
) 
{{ .Code }}
`

var openldapBuiltinAttributes = []attributetype{
	attributetype{
		Name: []string{"uidnumber"},
		Desc: "An integer uniquely identifying a user in an administrative domain",
		Equality: "integerMatch",
		Syntax: "1.3.6.1.4.1.1466.115.121.1.27",
		SingleValue: true},
	attributetype{
		Name: []string{"gidnumber"},
		Desc: "An integer uniquely identifying a group in an administrative domain",
		Equality: "integerMatch",
		Syntax: "1.3.6.1.4.1.1466.115.121.1.27",
		SingleValue: true}}

var pkg = flag.String("pkg", "objects", "go package to generate code for")
var builtin = flag.Bool("builtin", true, "include openldap builtin attribute definitions for uidnumber and gidnumber")
var objectsFile = flag.String("obj", "objects.json", "object description file")
var outputFile = flag.String("out", "objects.go", "output file")
var gofmt = flag.Bool("fmt", true, "call go fmt on the generated file")
var printOc = flag.Bool("objectclasses", false, "prints the parsed objectclasses")
var printAt = flag.Bool("attributes", false, "prints the parsed attributes")
var lexDebug = flag.Bool("lexdebug", false, "prints input cleaned of comments and each lexed token")
var schemaFiles []string

// replaces characters in names which are invalid names in go
var nameReplacer = strings.NewReplacer("-", "")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
  %v [options] schemafiles...
`, os.Args[0])
		flag.PrintDefaults()
		}

	flag.Parse()
	schemaFiles = flag.Args()
}

func readSchemas(schemaFiles []string) (string, error) {
	out := make([]string, len(schemaFiles))
	for i, fileName := range schemaFiles {
		bInput, err := ioutil.ReadFile(fileName)
		if err != nil {
			return "", err
		}
		out[i] = string(bInput)
	}

	return strings.Join(out, "\n"), nil
}

func main() {
	// Read the schema files and parse them. This fills the (global, because of yacc) variables
	// objectclassdefs and attributetypedefs with entrys.
	schemaInput, err := readSchemas(schemaFiles)
	if err != nil {
		log.Fatal(err)
	}
	l := newLexer(schemaInput, *lexDebug)
	yyParse(l)

	if *builtin {
		for _, v := range openldapBuiltinAttributes {
			attributetypedefs[v.Name[0]] = &v
		}
	}

	if *printAt {
		log.Println(attributetypedefs)
	}

	if *printOc {
		log.Println(objectclassdefs)
	}

	// Read the object definitions. Generate code for each one.
	objectsInput, err := ioutil.ReadFile(*objectsFile)
	if err != nil {
		log.Fatal(err)
	}
	objects, err := parseObjects(objectsInput)
	
	if err != nil {
		log.Fatal(err)
	}

	var code bytes.Buffer
	for _, v := range objects {
		codeString, err := v.Code()
		if err != nil {
			log.Fatal(err)
		}

		code.Write([]byte(codeString))
	}

	// Write the output
	data := struct {
		Package string
		Code    string
		SchemaFiles []string
		ObjectDesc string
		
	}{Package: *pkg, Code: code.String(), SchemaFiles: schemaFiles, ObjectDesc: string(objectsInput)}

	var output io.WriteCloser
	if *outputFile == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer output.Close()

	t := template.Must(template.New("code").Parse(temp))

	err = t.Execute(output, data)
	if err != nil {
		log.Fatal(err)
	}

	if *gofmt {
		cmd := exec.Command("go", "fmt", *outputFile)
		cmd.Run()
	}
}