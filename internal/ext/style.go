package ext

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

var boldKeywords = styles.Register(chroma.MustNewStyle("bold-keywords", chroma.StyleEntries{
	chroma.KeywordDeclaration: "bold",
	chroma.KeywordConstant:    "bold",
	chroma.KeywordNamespace:   "bold",
	chroma.KeywordType:        "bold",
	chroma.OperatorWord:       "bold",
}))
