package rconfig

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type Val []byte

func (h Val) FormatJson(v interface{}) error {
	return json.Unmarshal(h, v)
}

func (h Val) FormatYaml(v interface{}) error {
	return yaml.Unmarshal(h, v)
}

func (h Val) String() string {
	return string(h)
}

func (h Val) Bytes() []byte {
	return h
}
