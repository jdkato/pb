package platform

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jdkato/pb/internal/config"
)

type frontMatter struct {
	Title string
	Tags  []string
}

type post struct {
	body string
	meta frontMatter
}

// A Converter converts between platform-specific Markdown implementations.
type Converter struct {
	body []byte
}

func NewConverter(fp string) (Converter, error) {
	b, err := ioutil.ReadFile(fp)
	return Converter{body: b}, err
}

func (c Converter) To(platform string) error {
	switch strings.ToLower(platform) {
	case "medium":
		if config.Auth.Medium == "" {
			return errors.New("no medium token provided")
		}
		return c.toMedium()
	default:
		return fmt.Errorf("platform '%s' is not recognized", platform)
	}
}
