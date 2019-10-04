package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

func GetStructNames(file io.Reader) ([]string, error) {
	node, err := parser.ParseFile(token.NewFileSet(), "", file, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	visitor := &visitor{}
	ast.Walk(visitor, node)

	return visitor.StructNames, nil
}

type visitor struct {
	StructNames []string
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.TypeSpec:
		if _, ok := n.Type.(*ast.StructType); ok {
			v.StructNames = append(v.StructNames, n.Name.String())
		}
	}

	return v
}
