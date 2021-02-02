package config

import (
	"github.com/zdnscloud/cement/configure"
)

type ControllerConfig struct {
	Path        string `yaml:"-"`
	ETCDAddress string `yaml:"etcdAddress"`
	Interface   string `yaml:"interface"`
}

var gConf *ControllerConfig

func LoadConfig(path string) (*ControllerConfig, error) {
	var conf ControllerConfig
	conf.Path = path
	if err := conf.Reload(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *ControllerConfig) Reload() error {
	var newConf ControllerConfig
	if err := configure.Load(&newConf, c.Path); err != nil {
		return err
	}

	newConf.Path = c.Path
	*c = newConf
	gConf = &newConf
	return nil
}

func GetConfig() *ControllerConfig {
	return gConf
}
