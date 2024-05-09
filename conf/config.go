package conf

import (
	"big_market/common/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Database DatabaseConfig `yaml:"db"`
	Redis    RedisConfig    `yaml:"redis"`
	MQ       MQConfig       `yaml:"mq"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type RedisConfig struct {
	URL string `yaml:"url"`
}

type MQConfig struct {
	URL string `yaml:"url"`
}

var (
	Root *Config
)

func LoadConfig() Config {
	if Root == nil {
		env := os.Getenv("APP_ENV")
		dir, _ := os.Getwd()
		var err error
		if env == "" {
			log.Info("当前为dev环境")
			Root, err = readConfig(dir, "conf", "conf_dev.yml")
		} else if env == "prod" {
			log.Info("当前为prod环境")
			Root, err = readConfig(dir, "conf", "conf_prod.yml")
		}
		if err != nil {
			panic("error")
		}
	}
	return *Root
}

func readConfig(arg ...string) (*Config, error) {
	filename := filepath.Join(arg...)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
