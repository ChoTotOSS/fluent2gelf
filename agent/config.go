package agent

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type AgentConfig struct {
	Match     string
	Host      string
	Port      int
	Multiline bool
	Firstline *RecordMatch `yaml:"firstline_match"`
}

type RecordMatch struct {
	Begin   string
	IndexLT int `yaml:"index_lt"`
	Regexp  string
}

func NewConfig(match string, host string, port int, multiline bool) AgentConfig {
	return AgentConfig{
		Match:     match,
		Host:      host,
		Port:      port,
		Multiline: multiline,
	}
}

func LoadConfig(reader io.Reader) []AgentConfig {
	var configs []AgentConfig

	cfgContent, err := ioutil.ReadAll(reader)

	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(cfgContent, &configs); err != nil {
		panic(err)
	}
	return configs
}
