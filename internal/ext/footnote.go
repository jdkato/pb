package ext

import (
	"bytes"
	"strconv"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

var idx2ref = map[int]string{}

type footnoteHTMLRenderer struct {
	extension.FootnoteConfig
}

func newFootnoteHTMLRenderer(opts ...extension.FootnoteOption) renderer.NodeRenderer {
	r := &footnoteHTMLRenderer{
		FootnoteConfig: extension.NewFootnoteConfig(),
	}
	for _, opt := range opts {
		opt.SetFootnoteOption(&r.FootnoteConfig)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *footnoteHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(east.KindFootnoteLink, r.renderFootnoteLink)
	reg.Register(east.KindFootnoteBacklink, r.renderFootnoteBacklink)
	reg.Register(east.KindFootnote, r.renderFootnote)
	reg.Register(east.KindFootnoteList, r.renderFootnoteList)
}

func findClosestP(n gast.Node) []byte {
	p := n.Parent()
	for p.Kind() != gast.KindParagraph {
		p = p.Parent()
	}
	name, _ := n.Parent().AttributeString("id")
	return name.([]byte)
}

func (r *footnoteHTMLRenderer) renderFootnoteLink(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		n := node.(*east.FootnoteLink)
		is := strconv.Itoa(n.Index)
		name, _ := n.Parent().AttributeString("id")
		if name != nil {
			idx2ref[n.Index] = string(name.([]byte))
		}
		_, _ = w.WriteString(`<a href="#`)
		_, _ = w.Write(r.idPrefix(node))
		_, _ = w.WriteString(`fn:`)
		_, _ = w.WriteString(is)
		//_, _ = w.WriteString(`" data-name="`)
		//_, _ = w.Write(r.idPrefix(node))
		//_, _ = w.WriteString(`df:`)
		//_, _ = w.WriteString(is)
		_, _ = w.WriteString(`" class="`)
		_, _ = w.Write(applyFootnoteTemplate(r.FootnoteConfig.LinkClass,
			n.Index, n.RefCount))
		if len(r.FootnoteConfig.LinkTitle) > 0 {
			_, _ = w.WriteString(`" title="`)
			_, _ = w.Write(util.EscapeHTML(applyFootnoteTemplate(r.FootnoteConfig.LinkTitle, n.Index, n.RefCount)))
		}
		_, _ = w.WriteString(`" role="doc-noteref">`)

		_, _ = w.WriteString(formula("^" + is))
		_, _ = w.WriteString(`</a>`)
	}
	return gast.WalkContinue, nil
}

func (r *footnoteHTMLRenderer) renderFootnoteBacklink(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		n := node.(*east.FootnoteBacklink)
		is := strconv.Itoa(n.Index)
		_, _ = w.WriteString(`&#160;<a href="#`)
		if name, ok := idx2ref[n.Index]; ok {
			_, _ = w.WriteString(name)
		} else {
			_, _ = w.Write(r.idPrefix(node))
			_, _ = w.WriteString(`df:`)
			_, _ = w.WriteString(is)
		}
		_, _ = w.WriteString(`" class="`)
		_, _ = w.Write(applyFootnoteTemplate(r.FootnoteConfig.BacklinkClass, n.Index, n.RefCount))
		if len(r.FootnoteConfig.BacklinkTitle) > 0 {
			_, _ = w.WriteString(`" title="`)
			_, _ = w.Write(util.EscapeHTML(applyFootnoteTemplate(r.FootnoteConfig.BacklinkTitle, n.Index, n.RefCount)))
		}
		_, _ = w.WriteString(`" role="doc-backlink">`)
		_, _ = w.Write(applyFootnoteTemplate(r.FootnoteConfig.BacklinkHTML, n.Index, n.RefCount))
		_, _ = w.WriteString(`</a>`)
	}
	return gast.WalkContinue, nil
}

func (r *footnoteHTMLRenderer) renderFootnote(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*east.Footnote)
	is := strconv.Itoa(n.Index)
	if entering {
		_, _ = w.WriteString(`<li name="`)
		_, _ = w.Write(r.idPrefix(node))
		_, _ = w.WriteString(`fn:`)
		_, _ = w.WriteString(is)
		_, _ = w.WriteString(`" role="doc-endnote"`)
		if node.Attributes() != nil {
			html.RenderAttributes(w, node, html.ListItemAttributeFilter)
		}
		_, _ = w.WriteString(">\n")
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return gast.WalkContinue, nil
}

func (r *footnoteHTMLRenderer) renderFootnoteList(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	tag := "section"
	if r.Config.XHTML {
		tag = "div"
	}
	if entering {
		_, _ = w.WriteString("<")
		_, _ = w.WriteString(tag)
		_, _ = w.WriteString(` class="footnotes" role="doc-endnotes"`)
		if node.Attributes() != nil {
			html.RenderAttributes(w, node, html.GlobalAttributeFilter)
		}
		_ = w.WriteByte('>')
		if r.Config.XHTML {
			_, _ = w.WriteString("\n<hr />\n")
		} else {
			_, _ = w.WriteString("\n<hr>\n")
		}
		_, _ = w.WriteString("<ol>\n")
	} else {
		_, _ = w.WriteString("</ol>\n")
		_, _ = w.WriteString("</")
		_, _ = w.WriteString(tag)
		_, _ = w.WriteString(">\n")
	}
	return gast.WalkContinue, nil
}

func (r *footnoteHTMLRenderer) idPrefix(node gast.Node) []byte {
	if r.FootnoteConfig.IDPrefix != nil {
		return r.FootnoteConfig.IDPrefix
	}
	if r.FootnoteConfig.IDPrefixFunction != nil {
		return r.FootnoteConfig.IDPrefixFunction(node)
	}
	return []byte("")
}

func applyFootnoteTemplate(b []byte, index, refCount int) []byte {
	fast := true
	for i, c := range b {
		if i != 0 {
			if b[i-1] == '^' && c == '^' {
				fast = false
				break
			}
			if b[i-1] == '%' && c == '%' {
				fast = false
				break
			}
		}
	}
	if fast {
		return b
	}
	is := []byte(strconv.Itoa(index))
	rs := []byte(strconv.Itoa(refCount))
	ret := bytes.Replace(b, []byte("^^"), is, -1)
	return bytes.Replace(ret, []byte("%%"), rs, -1)
}

type footnote struct {
	options []extension.FootnoteOption
}

// Footnote is an extension that allow you to use PHP Markdown Extra Footnotes.
var Footnote = &footnote{
	options: []extension.FootnoteOption{},
}

// NewFootnote returns a new extension with given options.
func NewFootnote(opts ...extension.FootnoteOption) goldmark.Extender {
	return &footnote{
		options: opts,
	}
}

func (e *footnote) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(extension.NewFootnoteBlockParser(), 999),
		),
		parser.WithInlineParsers(
			util.Prioritized(extension.NewFootnoteParser(), 101),
		),
		parser.WithASTTransformers(
			util.Prioritized(extension.NewFootnoteASTTransformer(), 999),
		),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(newFootnoteHTMLRenderer(e.options...), 500),
	))
}
