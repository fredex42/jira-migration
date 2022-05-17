package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ScriptKey struct {
	User string `yaml:"user"`
	Key  string `yaml:"key"`
}

func LoadScriptKey(path *string) (*ScriptKey, error) {
	f, err := os.Open(*path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var key ScriptKey
	err = yaml.Unmarshal(content, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}
