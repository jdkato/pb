package ext

import (
	_ "image/png"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	grh "github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type mathHTMLRenderer struct {
	grh.Config

	usePNG bool
}

func newMathHTMLRenderer(png bool, opts ...grh.Option) renderer.NodeRenderer {
	r := &mathHTMLRenderer{
		Config: grh.NewConfig(),
		usePNG: png}

	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}

	return r
}

func (r *mathHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindMath, r.render)
}

func (r *mathHTMLRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	var err error

	if entering {
		n := node.(*mathNode)
		if n.IsInline {
			w.WriteString("<em>" + formula(string(n.Value)) + "</em>")
		} else {
			err = toImg(w, n, r.usePNG)
		}

		if err != nil {
			panic(err) // Shouldn't happen
		}
	}

	return ast.WalkContinue, nil
}
