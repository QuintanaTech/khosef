package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"path"
	"strings"
)

var fileName = ".khosef.toml"

type Config struct {
	Provider   string   `toml:"provider"`
	Secrets    []string `toml:"secrets"`
	contextDir string
}

func ReadConfig(contextDir string) (*Config, error) {
	c := NewConfig(contextDir)
	d := os.DirFS(c.contextDir)
	if _, err := toml.DecodeFS(d, fileName, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) Save() error {
	f, err := os.OpenFile(path.Join(c.contextDir, fileName), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(c)
}

func (c *Config) GetSecretDefinitions() []*SecretDefinition {
	s := make([]*SecretDefinition, len(c.Secrets))
	for i, text := range c.Secrets {
		parsed := strings.Split(text, ":")
		s[i] = &SecretDefinition{
			OutputPath: parsed[0],
			SecretId:   parsed[1],
		}
	}

	return s
}

func (c *Config) GetContextDir() string {
	return c.contextDir
}

func (c *Config) FindFile(filePath string) (int, *SecretDefinition) {
	for i, s := range c.GetSecretDefinitions() {
		p := path.Join(c.contextDir, filePath)

		if p == path.Join(c.contextDir, s.OutputPath) {
			return i, s
		}
	}

	return -1, nil
}

func NewConfig(contextDir string) *Config {
	return &Config{contextDir: contextDir}
}
