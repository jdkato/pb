package ext

import (
	"bytes"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var contextKey = parser.NewContextKey()

type data struct {
	Equation    []byte
	Node        gast.Node
	Prev        gast.Node
	Transformed bool
}

type mathBlockParser struct{}

var defaultMathBParser = &mathBlockParser{}

func newMathBlockParser() parser.BlockParser {
	return defaultMathBParser
}

func isSeparator(line []byte, sep byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	for i := 0; i < len(line); i++ {
		if line[i] != sep {
			return false
		}
	}
	return true
}

func (b *mathBlockParser) Trigger() []byte {
	return []byte{'$'}
}

func (b *mathBlockParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if isSeparator(line, '$') {
		return gast.NewTextBlock(), parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *mathBlockParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if isSeparator(line, '$') && !util.IsBlank(line) {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *mathBlockParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	var buf bytes.Buffer

	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}
	d := &data{Node: node, Equation: buf.Bytes(), Prev: node.PreviousSibling()}

	escapedList := pc.ComputeIfAbsent(contextKey,
		func() interface{} {
			return []*data{}
		}).([]*data)
	escapedList = append(escapedList, d)

	pc.Set(contextKey, escapedList)
	node.Parent().RemoveChild(node.Parent(), node)
}

func (b *mathBlockParser) CanInterruptParagraph() bool {
	return false
}

func (b *mathBlockParser) CanAcceptIndentedLine() bool {
	return true
}

type astTransformer struct{}

var defaultASTTransformer = &astTransformer{}

func (a *astTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	lst := pc.Get(contextKey)
	if lst == nil {
		return
	}

	pc.Set(contextKey, nil)
	for _, d := range lst.([]*data) {
		if d.Transformed {
			continue
		}

		img := newMathNode(false, d.Equation)
		node.InsertAfter(node, d.Prev, img)

		d.Transformed = true
	}
}

type math2Img struct {
	ConvertToPNG bool
}

// MathOption is a functional option type for this extension.
type MathOption func(*math2Img)

// WithConvertToPNG is a functional option that instructs Math2Image to perform
// a local conversion from SVG to PNG using InkScape.
//
// This is required for certain platforms -- such as Medium -- that don't
// support SVGs.
func WithConvertToPNG() MathOption {
	return func(m *math2Img) {
		m.ConvertToPNG = true
	}
}

// New returns a new Math2Img extension.
func NewMath2Img(opts ...MathOption) goldmark.Extender {
	e := &math2Img{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Math2Image is a extension for Goldmark that converts block-style
// (`$$ ... $$`) math equations into an image sutiable for platforms without
// math typesetting support.
var Math2Img = &math2Img{}

func (e *math2Img) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(newMathBlockParser(), 0),
		),
		parser.WithASTTransformers(
			util.Prioritized(defaultASTTransformer, 0),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newMathHTMLRenderer(e.ConvertToPNG), 0),
	))
}
