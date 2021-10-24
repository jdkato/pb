package ext

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jdkato/pb/internal/cli"
	"github.com/jdkato/pb/internal/config"
	msdk "github.com/medium/medium-sdk-go"
	"github.com/pterm/pterm"
	"github.com/yuin/goldmark/util"
)

// FromLocalToMedium uploads a local file to Medium.
func FromLocalToMedium(src string) string {
	var contentType string

	if strings.Contains(src, "http") {
		return src
	}

	src = filepath.Join(cli.Flags.ImageDir, src)
	switch strings.ToLower(filepath.Ext(src)) {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".tiff":
		contentType = "image/tiff"
	default:
		panic(fmt.Sprintf("unsupported image type '%s'", src))
	}

	m := msdk.NewClientWithAccessToken(config.Auth.Medium)

	ret, err := m.UploadImage(msdk.UploadOptions{
		FilePath:    src,
		ContentType: contentType,
	})

	if err != nil {
		panic(err)
	}

	pterm.Success.Printf(
		"Uploaded image '%s' (%s)\n", filepath.Base(src), contentType)

	return ret.URL
}

func figure(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func toPNG(w util.BufWriter, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file1, err := ioutil.TempFile("", "img.*.svg")
	if err != nil {
		return err
	}
	defer os.Remove(file1.Name())

	_, err = io.Copy(file1, resp.Body)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"inkscape",
		"-h",
		"64",
		file1.Name(),
		"-o",
		"temp.png")

	if err := cmd.Run(); err != nil {
		return err
	}

	file2, err := os.Open("temp.png")
	if err != nil {
		return err
	}
	defer os.Remove("temp.png")

	image, _, err := image.DecodeConfig(file2)
	if err != nil {
		return err
	}
	m := msdk.NewClientWithAccessToken(config.Auth.Medium)

	ret, err := m.UploadImage(msdk.UploadOptions{
		FilePath:    "temp.png",
		ContentType: "image/png",
	})
	if err != nil {
		return err
	}

	img := figure(`<img data-width="{w}" data-height="{h}" src="{src}">`,
		"{w}", strconv.Itoa(image.Width),
		"{h}", strconv.Itoa(image.Height),
		"{src}", ret.URL)

	w.WriteString(img)
	return nil
}

func toImg(w util.BufWriter, n *mathNode, png bool) error {
	value := url.QueryEscape(string(n.Value))

	page := "https://math.now.sh?color=black&from=" + value
	if png {
		return toPNG(w, page)
	}

	w.WriteString(fmt.Sprintf(`<img src="%s">`, page))
	return nil
}
