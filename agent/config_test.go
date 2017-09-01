package agent

import (
	"strings"
	"testing"
)

var dumpConfig = `
- match: docker.*.abc
  host: localhost
  port: 12701
  multiline: true
  begin_record: 201
`

func TestReadConfig(t *testing.T) {
	reader := strings.NewReader(dumpConfig)
	configs := LoadConfig(reader)

	if configs[0].Match != "docker.*.abc" {
		t.Fail()
	}

	if configs[0].Host != "localhost" {
		t.Fail()
	}

	if configs[0].Port != 12701 {
		t.Fail()
	}
}
