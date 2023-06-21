package config

type Config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Redis struct {
		Port     string `yaml:"port" envconfig:"REDIS_PORT"`
		Host     string `yaml:"host" envconfig:"REDIS_HOST"`
		DB       int    `yaml:"db" envconfig:"REDIS_DB"`
		Password string `yaml:"password" envconfig:"REDIS_PASSWORD"`
	} `yaml:"redis"`
	FileName string `yaml:"fileName" envconfig:"FILE_NAME"`
}
