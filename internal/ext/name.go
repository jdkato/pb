package ext

import (
	"fmt"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type namedPTransformer struct{}

var defaultNamedPTransformer = &namedPTransformer{}
var index = 0

func (g *namedPTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	_ = gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if n.Kind() == gast.KindParagraph || n.Kind() == gast.KindListItem {
			if entering {
				index++
				name := fmt.Sprintf("%04d", index)
				n.SetAttributeString("id", []byte(name))
			}
		}
		return gast.WalkContinue, nil
	})
}

type namedP struct{}

// NamedP is an extension that assigns each paragraph a unique `name=` ID.
var NamedP = &namedP{}

func (e *namedP) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(defaultNamedPTransformer, 0),
		),
	)
}
