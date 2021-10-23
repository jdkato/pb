package platform

import (
	"bytes"
	"regexp"

	"github.com/mitchellh/mapstructure"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	grh "github.com/yuin/goldmark/renderer/html"

	"github.com/jdkato/pb/internal/ext"
	highlighting "github.com/yuin/goldmark-highlighting"
)

// mediumMd is an extension designed to accomdate markup limitations of
// https://medium.com/.
var mediumMd = goldmark.New(
	goldmark.WithExtensions(
		// Assigns `name=` IDs to each paragraph so that we have something for
		// backlinks to reference.
		ext.NamedP,

		// Address Medium's lack of math typesetting: inline equations (`$...$`)
		// are replaced with unicode symbols, while block equations are replaced
		// with an image.
		ext.NewMath2Img(ext.WithConvertToPNG()),
		ext.Math2Unicode,

		// Add support for [^1]-style footnotes.
		ext.Footnote,

		// Add support for YAML-based front matter.
		//
		// This is used to add platform-specifc metadata like a title, summary, and
		// tags.
		meta.Meta,

		// Address Medium's lack of table support: tables are converted into an
		// ASCII-formatted table inside a code block.
		ext.AsciiTable,

		// Address Medium's lack of syntax-highlighting by bolding keywords.
		highlighting.NewHighlighting(highlighting.WithStyle("bold-keywords")),
	),
	goldmark.WithRendererOptions(
		grh.WithUnsafe(),
	),
)

var style2Tag = regexp.MustCompile(`<span style="font-weight:bold">(.+)</span>`)
var p2Name = regexp.MustCompile(`<p id="(p-\d)"`)

func toMediumMarkdown(b []byte) (post, error) {
	var buf bytes.Buffer

	ctx := parser.NewContext()
	p := post{}

	err := mediumMd.Convert(b, &buf, parser.WithContext(ctx))
	if err != nil {
		return p, err
	}

	err = mapstructure.Decode(meta.Get(ctx), &p.meta)
	if err != nil {
		return p, err
	}

	html := style2Tag.ReplaceAllString(buf.String(), "<strong>${1}</strong>")
	html = p2Name.ReplaceAllString(html, `<p name="${1}"`)

	p.body = html
	return p, nil
}
