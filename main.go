package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"verve_testcase/pkg/service"
	"verve_testcase/pkg/storage"
)

func main() {
	fileName := "data/promotions.csv"
	store := storage.CreateStorage()
	srv := service.CreateService(store)
	err := srv.Reload(fileName)
	if err != nil {
		fmt.Println("Service init failed")
	}

	router := gin.Default()
	router.GET("/promotions/:id", srv.GetPromotionById)
	router.POST("/promotions/reload", srv.ReloadPromotions)

	err = router.Run("localhost:1321")
	if err != nil {
		fmt.Println("Server run error", err)
		return
	}
}
