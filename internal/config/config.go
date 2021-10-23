package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/adrg/xdg"
	"github.com/cristalhq/aconfig"
)

var Auth *Config

func init() {
	Auth = new()
}

type Config struct {
	Medium   string `json:"medium"`   // The user-generated Medium Integration token
	DEV      string `json:"dev"`      // The user-generated DEV Community API key
	Hashnode string `json:"hashnode"` // The user-generated Hashnode Personal Access Token
	GitHub   string `json:"github"`   // The user-generated GitHub Person Access Token

	Path string `json:"path"`
}

func new() *Config {
	fp, err := xdg.ConfigFile("pb/auth.json")
	if err != nil {
		panic(err)
	}
	cfg := &Config{Path: fp}

	loader := aconfig.LoaderFor(cfg, aconfig.Config{
		SkipDefaults:       true,
		AllFieldRequired:   false,
		SkipFlags:          true,
		EnvPrefix:          "PB",
		AllowUnknownEnvs:   true,
		AllowUnknownFields: true,
		Files:              []string{fp},
	})

	if err = loader.Load(); err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) Save() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Path, b, 0644)
}

func (c *Config) AddToken(token, platform string) error {
	switch platform {
	case "Medium":
		c.Medium = token
	case "DEV":
		c.DEV = token
	case "Hashnode":
		c.Hashnode = token
	default:
		return fmt.Errorf("unknown platform '%s'", platform)
	}
	return nil
}
