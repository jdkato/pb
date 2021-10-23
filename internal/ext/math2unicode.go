package ext

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var latexInlineRegexp = regexp.MustCompile(`^\$[^$]*\$`)

type mathInlineParser struct {
}

var defaultMathIParser = &mathInlineParser{}

// newMathIParser return a new InlineParser that parses
// latex expressions.
func newMathIParser() parser.InlineParser {
	return defaultMathIParser
}

func (s *mathInlineParser) Trigger() []byte {
	return []byte{'$'}
}

func (s *mathInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	m := latexInlineRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	block.Advance(m[1])

	value := block.Value(text.NewSegment(
		segment.Start+1,
		segment.Start+m[1]-1))

	return newMathNode(true, value)
}

func (s *mathInlineParser) CloseBlock(parent ast.Node, pc parser.Context) {}

type math2Unicode struct{}

// Math2Unicode is a extension for Goldmark that converts inline-style
// (`$ ... $`) math equations into their unicode-equivalent characters for use
// on platform without math typesetting support.
var Math2Unicode = &math2Unicode{}

func (e *math2Unicode) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(newMathIParser(), 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newMathHTMLRenderer(false), 500),
	))
}
