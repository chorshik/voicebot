package system

import (
	"encoding/json"
	"io/ioutil"
	conf"github.com/ebladrocher/voicebot/system/model"
)

// Path ...
const (
	path = "conf/config.json"
)

// ReadConfig ...
func ReadConfig() (cfg *conf.Config, err error) {

	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg = &conf.Config{}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return
}
