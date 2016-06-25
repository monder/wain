package wain

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigUrl struct {
	Pattern  string
	Original struct {
		Bucket string
		Key    string
	}
	Cache *struct {
		Bucket string
		Key    string
	}
}

type Config struct {
	Port    int
	Buckets []struct {
		Name         string
		Region       string
		AccessKey    string `yaml:"accessKey"`
		AccessSecret string `yaml:"accessSecret"`
	}
	Urls []ConfigUrl
}

func ReadConfig(path string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	if config.Port == 0 {
		config.Port = 3000
	}
	return &config, nil
}
