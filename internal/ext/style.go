package ext

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

var boldKeywords = styles.Register(chroma.MustNewStyle("bold-keywords", chroma.StyleEntries{
	chroma.KeywordDeclaration: "bold",
	chroma.KeywordConstant:    "bold",
	chroma.KeywordNamespace:   "bold",
	chroma.KeywordReserved:    "bold",
	chroma.Keyword:            "bold",
	chroma.OperatorWord:       "bold",
}))
