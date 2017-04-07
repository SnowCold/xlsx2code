package base

import (
	"encoding/json"
	"io/ioutil"
)

const (
	CommentLine   = 0
	BelongLine    = 1
	TypeLine      = 2
	FiledNameLine = 3
	IgnoreLine    = 4
	Delimiter     = "\t"
)

type Config struct {
	DontOutputCode []string
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) Load(jsonFilePath string) {
	confByte, _ := ioutil.ReadFile(jsonFilePath)
	json.Unmarshal(confByte, config)
}

func (config *Config) CheckNeedOutputCode(fileName string) bool {
	for i := 0; i < len(config.DontOutputCode); i++ {
		if config.DontOutputCode[i] == fileName {
			return false
		}
	}
	return true
}
