package ext

import (
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var defaultTableParser = &tableBlockParser{}
var reIsTable = regexp.MustCompile(`\|.+\|`)

type tableBlockParser struct{}

func newTableBlockParser() parser.BlockParser {
	return defaultTableParser
}

func (b *tableBlockParser) Trigger() []byte {
	return []byte{'|'}
}

func (b *tableBlockParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, segment := reader.PeekLine()
	if reIsTable.Match(line) {
		n := gast.NewCodeBlock()
		n.Lines().Append(segment)
		return n, parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *tableBlockParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if util.IsBlank(line) {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *tableBlockParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {}

func (b *tableBlockParser) CanInterruptParagraph() bool {
	return false
}

func (b *tableBlockParser) CanAcceptIndentedLine() bool {
	return false
}

type asciiTable struct{}

// AsciiTable is a extension for Goldmark that wraps Markdown tables in code
// blocks.
var AsciiTable = &asciiTable{}

func (e *asciiTable) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(newTableBlockParser(), 0),
		),
	)
}
