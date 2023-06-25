package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"verve_challenge/pkg/config"
	"verve_challenge/pkg/service"
	"verve_challenge/pkg/storage"
)

func main() {
	var cfg config.Config
	readFile(&cfg)
	readEnv(&cfg)

	store, err := storage.CreateStorage(cfg)
	if err != nil {
		processError(err)
	}
	srv := service.CreateService(store)

	err = srv.Reload(cfg.FileName)
	if err != nil {
		processError(err)
	}

	router := gin.Default()
	router.GET("/promotions/:id", srv.GetPromotionById)
	router.POST("/promotions/reload", srv.ReloadPromotions)

	err = router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *config.Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *config.Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
