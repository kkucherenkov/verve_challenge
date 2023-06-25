package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"verve_challenge/pkg/file_processor"
	"verve_challenge/pkg/storage"
)

type reloadRequest struct {
	Filename string `json:"filename"`
}
type Service struct {
	storage *storage.Storage
}

func (srv *Service) Reload(fileName string) error {
	srv.storage.Clean()
	err := file_processor.ProcessFile(fileName, srv.storage)

	if err != nil {
		fmt.Println("File processing failed")
		return err
	}
	return nil
}

func (srv *Service) ReloadPromotions(context *gin.Context) {
	var requestBody reloadRequest
	if err := context.BindJSON(&requestBody); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Wrong request format"})
		return
	}
	err := srv.Reload(requestBody.Filename)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "File parse error"})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func (srv *Service) getById(id string) *storage.Item {
	return srv.storage.Get(id)
}

func (srv *Service) GetPromotionById(context *gin.Context) {
	id := context.Param("id")
	item := srv.getById(id)
	if item != nil {
		context.IndentedJSON(http.StatusOK, item)
	} else {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "promotion not found"})
	}
}

func CreateService(storage *storage.Storage) *Service {
	return &Service{storage: storage}
}
