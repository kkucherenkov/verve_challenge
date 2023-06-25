package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Storage struct {
		Port        string `yaml:"port" envconfig:"STORAGE_PORT"`
		Host        string `yaml:"host" envconfig:"STORAGE_HOST"`
		MetricsPort string `yaml:"metricsPort" envconfig:"STORAGE_METRICS_PORT"`
	} `yaml:"storage"`
	WebApi struct {
		Port        string `yaml:"port" envconfig:"WEB_API_PORT"`
		Host        string `yaml:"host" envconfig:"WEB_API_HOST"`
		StorageAddr string `yaml:"clientAddr" envconfig:"STORAGE_ADDR"`
		MetricsPort string `yaml:"metricsPort" envconfig:"WEB_API_METRICS_PORT"`
	} `yaml:"webApi"`
	FileName string `yaml:"fileName" envconfig:"FILE_NAME"`
}

func ReadFile(cfg *Config) error {
	f, err := os.Open("config.yml")
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func ReadEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}
	return nil
}
