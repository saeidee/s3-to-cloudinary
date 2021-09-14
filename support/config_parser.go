package support

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type configParser struct{}

type Config struct {
	Version string   `yaml:"version"`
	Buckets []string `yaml:"buckets"`
	MaxKeys int64    `yaml:"maxKeys"`
	Secrets struct {
		Cloudinary struct {
			CloudName    string `yaml:"cloudName"`
			ApiKey       string `yaml:"apiKey"`
			ApiSecret    string `yaml:"apiSecret"`
			UploadPrefix string `yaml:"uploadPrefix"`
		}
		AWS struct {
			AccessKeyID     string `yaml:"accessKeyId"`
			SecretAccessKey string `yaml:"secretAccessKey"`
		}
	}
}

func NewConfigParser() *configParser {
	return &configParser{}
}

func (cp *configParser) Parse(filename string) (*Config, error) {
	var cfg Config

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(source, &cfg)

	return &cfg, err
}

func (c *Config) HasUploadPrefix() bool {
	return len(c.Secrets.Cloudinary.UploadPrefix) != 0
}
