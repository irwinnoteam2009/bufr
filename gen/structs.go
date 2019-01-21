package main

import (
	"fmt"
	"go/ast"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

var (
	tpl = template.Must(template.New("field").Parse(`
	{{if .Type eq "bool"}}
	if v, err := reader.Read{{if .Op eq "bits"}}Bit{{else}}Byte{{end}}(); err != nil {
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

	tplArr = template.Must(template.New("array").Parse(`
	// read {{.Name}}
	for i := 0; i < {{.Wants}}; i++ {
		v, err := reader.Read {{ if .Op eq "bits" }}Bit{{else}}Byte{{end}}()
		if err != nil {
			return err
		}
		t.{{.Name}}[i] = v
	}
	`))
)

type fieldInfo struct {
	Name  string
	Type  string
	Op    string
	Wants string
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
	fmt.Println("doc:", t.Doc.Text(), "comments:", t.Comment.Text())
	fmt.Println(t.Name.Name)

	for _, f := range s.Fields.List {
		var vtag string
		if f.Tag != nil {
			tag := reflect.StructTag(f.Tag.Value[1 : len(f.Tag.Value)-1])
			vtag = tag.Get("bufr")
		}
		processField(f, vtag)
		// fmt.Println()
	}

	return nil
}

const maxCount = 111

// Test is lol
type Test struct {
	Magic   [4]byte `bufr:"bytes:4"`       // read 4 bytes
	IsBool  bool    `bufr:"bits:1,skip:7"` // read 1 bit as bool
	IsBool2 bool    `bufr:"bytes:1"`       // read byte as bool
	IsInt   int     `bufr:"bytes:3"`       // read 3 bytes as int
	IsInt2  int     // read 4 bytes
	Skip    bool    `bufr:"-"` // skip
	IsByte  byte    // read 1 byte
	IsFloat float64 `bufr:"code:005002"`
	IsArr   []int
	IsArr2  [7]bool
	IsArr3  [maxCount]int
	Map     map[string]int
	Pointer *int
}

func processField(field *ast.Field, tag string) {
	fmt.Printf("%s %T\n", field.Names[0].Name, field.Type)

	if tag == "-" {
		return
	}
	values := getTagValues(tag)
	info := fieldInfo{
		Name: field.Names[0].Name,
		Op:   "bytes",
	}

	// by type
	var typLen string
	switch t := field.Type.(type) {
	case *ast.ArrayType:
		if t.Len != nil {
			switch l := t.Len.(type) {
			case *ast.BasicLit:
				typLen = l.Value
			case *ast.Ident:
				typLen = l.Name
			}
		}
		// fmt.Printf("%T\n", l)
	case *ast.MapType, *ast.StarExpr:
		// pointers
		fmt.Println("SKIP UNUSED")
		return
	case *ast.Ident:
		// simple types
		typLen = getSizeByType(t.Name)
		fmt.Println("simple type: ", t.Name)
	}
	info.Wants = typLen

	// override info
	for k, v := range values {
		switch k {
		case "bits", "bytes":
			info.Op = k
			_, err := strconv.Atoi(v)
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

	fmt.Printf("%+v\n", info)
	tplArr.Execute(os.Stdout, info)
}

func getSizeByType(typ string) string {
	size := map[string]string{
		"bool":    "1",
		"byte":    "1",
		"int8":    "1",
		"int16":   "2",
		"int32":   "4",
		"int64":   "8", // 4
		"uint8":   "1",
		"uint16":  "2",
		"uint32":  "4",
		"uint64":  "8", // 4
		"int":     "4",
		"float32": "4",
		"float64": "8", // 4
	}
	return size[typ]
}
