package agent

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type AgentConfig struct {
	Match       string
	Host        string
	Port        int
	Multiline   bool
	BeginRecord string `yaml:"begin_record"`
}

func NewConfig(match string, host string, port int, multiline bool, begin string) AgentConfig {
	return AgentConfig{
		Match:       match,
		Host:        host,
		Port:        port,
		Multiline:   multiline,
		BeginRecord: begin,
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
