package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// out, err := os.Create(os.Args[2])
	// if err != nil {
	// 	panic(err)
	// }

	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			fmt.Println(t.Comment.Text(), t.Doc.Text())
			if currStruct, ok := t.Type.(*ast.StructType); ok {
				if err := processStruct(t, currStruct); err != nil {
					panic(err)
				}
			}
		}

		return true
	})
}
