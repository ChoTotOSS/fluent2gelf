package agent

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Match     string
	Host      string
	Port      int
	Multiline *FirstlineMatch
}

type FirstlineMatch struct {
	Begin   string
	IndexLT int `yaml:"index_lt"`
	Regexp  string
}

func NewConfig(match string, host string, port int, multiline bool) Config {
	return Config{
		Match:     match,
		Host:      host,
		Port:      port,
		Multiline: nil,
	}
}

func LoadConfig(reader io.Reader) []Config {
	var configs []Config

	cfgContent, err := ioutil.ReadAll(reader)

	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(cfgContent, &configs); err != nil {
		panic(err)
	}
	return configs
}
