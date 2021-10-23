package ext

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

var kindMath = ast.NewNodeKind("Math")

type mathNode struct {
	ast.BaseInline
	IsInline bool
	Value    []byte
}

// Dump implements Node.Dump.
func (n *mathNode) Dump(source []byte, level int) {
	m := map[string]string{
		"Inline": fmt.Sprintf("%v", n.IsInline),
		"Vaule":  fmt.Sprintf("%v", n.Value),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

// Kind implements Node.Kind.
func (n *mathNode) Kind() ast.NodeKind {
	return kindMath
}

func newMathNode(isInline bool, value []byte) *mathNode {
	return &mathNode{
		IsInline: isInline,
		Value:    value,
	}
}

var kindASCII = ast.NewNodeKind("ASCII")

type asciiNode struct {
	ast.BaseBlock
	Value []byte
}

// Dump implements Node.Dump.
func (n *asciiNode) Dump(source []byte, level int) {
	m := map[string]string{
		"Vaule": fmt.Sprintf("%v", n.Value),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

// Kind implements Node.Kind.
func (n *asciiNode) Kind() ast.NodeKind {
	return kindASCII
}

func newASCIINode(value []byte) *asciiNode {
	return &asciiNode{Value: value}
}
