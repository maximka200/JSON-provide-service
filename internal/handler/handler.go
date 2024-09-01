package handler

import (
	"jps/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Storage storage.Storage
}

func NewHandler(strg storage.Storage) Handler {
	return Handler{Storage: strg}
}

func (h *Handler) InitRouter() *gin.Engine {
	engine := gin.New()
	service := engine.Group("/service")
	{
		service.POST("/newJSON", h.NewJSON)
		service.GET("/getJSON", h.GetJSON)
		service.DELETE("/deleteJSON", h.DeleteJSON)
	}
	return engine
}
