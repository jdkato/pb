package ext

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/yuin/goldmark"
)

func Test2Unicode(t *testing.T) {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(
			Math2Unicode,
		),
	)
	dir := filepath.Join(testdata, "inline")

	cases, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, c := range cases {
		exp, _ := ioutil.ReadFile(filepath.Join(dir, c.Name(), "test.md"))
		obs, _ := ioutil.ReadFile(filepath.Join(dir, c.Name(), "test.html"))

		md.Convert(exp, &buf)
		if !bytes.Equal(buf.Bytes(), obs) {
			t.Fatalf("%s failed", string(buf.String()))
		}

		buf.Reset()
	}
}
