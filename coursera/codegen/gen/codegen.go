// go build gen/* && ./codegen.exe pack/unpack.go  pack/marshaller.go
// go run pack/*
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

type tpl struct {
	FieldName string
}

var (
	intTpl = template.Must(template.New("intTpl").Parse(`
	// {{.FieldName}}
	var {{.FieldName}}Raw uint32
	binary.Read(r, binary.LittleEndian, &{{.FieldName}}Raw)
	in.{{.FieldName}} = int({{.FieldName}}Raw)
`))

	strTpl = template.Must(template.New("strTpl").Parse(`
	// {{.FieldName}}
	var {{.FieldName}}LenRaw uint32
	binary.Read(r, binary.LittleEndian, &{{.FieldName}}LenRaw)
	{{.FieldName}}Raw := make([]byte, {{.FieldName}}LenRaw)
	binary.Read(r, binary.LittleEndian, &{{.FieldName}}Raw)
	in.{{.FieldName}} = string({{.FieldName}}Raw)
`))
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import "encoding/binary"`)
	fmt.Fprintln(out, `import "bytes"`)
	fmt.Fprintln(out) // empty line

	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			fmt.Printf("SKIP %#T is not *ast.GenDecl\n", f)
			continue
		}
	SPECS_LOOP:
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				fmt.Printf("SKIP %#T is not ast.TypeSpec\n", spec)
				continue
			}

			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				fmt.Printf("SKIP %#T is not ast.StructType\n", currStruct)
				continue
			}

			if g.Doc == nil {
				fmt.Printf("SKIP struct %#v doesnt have comments\n", currType.Name.Name)
				continue
			}

			needCodegen := false
			for _, comment := range g.Doc.List {
				needCodegen = needCodegen || strings.HasPrefix(comment.Text, "// cgen: binpack")
			}
			if !needCodegen {
				fmt.Printf("SKIP struct %#v doesnt have cgen mark\n", currType.Name.Name)
				continue SPECS_LOOP
			}

			fmt.Printf("process struct %s\n", currType.Name.Name)
			fmt.Printf("\tgenerating Unpack method\n")

			fmt.Fprintln(out, "func (in *"+currType.Name.Name+") Unpack(data []byte) error {")
			fmt.Fprintln(out, "	r := bytes.NewReader(data)")

		FIELDS_LOOP:
			for _, field := range currStruct.Fields.List {

				if field.Tag != nil {
					tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
					if tag.Get("cgen") == "-" {
						continue FIELDS_LOOP
					}
				}

				fieldName := field.Names[0].Name
				fileType := field.Type.(*ast.Ident).Name

				fmt.Printf("\tgenerating code for field %s.%s\n", currType.Name.Name, fieldName)

				switch fileType {
				case "int":
					intTpl.Execute(out, tpl{fieldName})
				case "string":
					strTpl.Execute(out, tpl{fieldName})
				default:
					log.Fatalln("unsupported", fileType)
				}
			}

			fmt.Fprintln(out, "	return nil")
			fmt.Fprintln(out, "}") // end of Unpack func
			fmt.Fprintln(out)      // empty line

		}
	}
}

// go build gen/* && ./codegen.exe pack/unpack.go  pack/marshaller.go
// go run pack/*
