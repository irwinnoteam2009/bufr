package main

import (
	"fmt"
	"go/ast"
	"html/template"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	tpl = template.Must(template.New("field").Parse(`
	{{if .Type eq "bool"}}
	if v, err := reader.Read{{if .Op eq "bits"}}Bit{{else if .Op eq "bytes"}}Byte{{end}}(); err != nil {
		return err
	}
	t.{{.Name}} = v != 0
	{{else if .Type eq "int"}}

	{{else if .Type eq "byte"}}

	{{end}}

	{{if .Skip}}
	if _, err := reader.ReadBits({{.Wants}} {{if .Op eq "bytes"}} * 8{{end}}); err != nil {
		return err
	}
	{{end}}
	`))
)

type fieldInfo struct {
	Name  string
	Type  string
	Op    string
	Wants int
	Skip  int
}

func getTagValues(tag string) map[string]string {
	values := make(map[string]string)
	// split tag value to kv by , and =
	arr := strings.Split(tag, ",")
	for _, value := range arr {
		arr := strings.Split(value, ":")
		if len(arr) == 1 {
			values[arr[0]] = ""
		} else {
			values[arr[0]] = arr[1]
		}
	}
	return values
}

func processStruct(t *ast.TypeSpec, s *ast.StructType) error {
	fmt.Println(t.Doc.Text(), t.Comment.Text())
	if comment := t.Doc; comment != nil {
		fmt.Println(comment.Text())
	}
	fmt.Println(t.Name.Name)

	for _, f := range s.Fields.List {
		// fmt.Printf("%v ", f)
		if f.Tag != nil {
			tag := reflect.StructTag(f.Tag.Value[1 : len(f.Tag.Value)-1])
			v := tag.Get("bufr")
			processField(f, v)
		}
		// fmt.Println()
	}

	return nil
}

type test struct {
	Magic   [4]byte `bufr:"bytes:4"`       // read 4 bytes
	IsBool  bool    `bufr:"bits:1,skip:7"` // read 1 bit as bool
	IsBool2 bool    `bufr:"bytes:1"`       // read byte as bool
	IsInt   int     `bufr:"bytes:3"`       // read 3 bytes as int
	IsInt2  int     // read 4 bytes
	Skip    bool    `bufr:"-"` // skip
	IsByte  byte    // read 1 byte

}

func processField(field *ast.Field, tag string) {
	if tag == "-" {
		return
	}
	values := getTagValues(tag)
	info := fieldInfo{
		Name: field.Names[0].Name,
	}

	for k, v := range values {
		switch k {
		case "bits", "bytes":
			info.Op = k
			v, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			info.Wants = v
		case "skip":
			v, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			info.Skip = v
		case "type":

		}
	}
	// TODO: добавить обработку по типу

	// TODO: добавить Op и Wants в зависимости от типа
	if info.Op == "" {

	}

	fmt.Printf("%+v\n", info)
	tpl.Execute(os.Stdout, info)
}
