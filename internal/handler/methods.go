package handler

import (
	"errors"
	"fmt"
	"io"
	storageErr "jps/internal/storage/postgresql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Id struct {
	Id int `json:"id"`
}

// слой десериализации данных из http req

func (h *Handler) NewJSON(c *gin.Context) {
	json, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("cannot read message body: %s", err))
		return
	}
	defer c.Request.Body.Close()

	jsonStr := string(json)

	id, err := h.Storage.NewJSON(jsonStr)
	if err != nil {
		if errors.Is(err, storageErr.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, "invalid credentials")
			return
		}
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("cannot get json id: %s", err))
		return
	}

	c.JSON(http.StatusOK, map[string]int{"id": id})
}

func (h *Handler) DeleteJSON(c *gin.Context) {
	var id Id

	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("not correct request: %s", err))
		return
	}

	if err := h.Storage.DeleteJSON(id.Id); err != nil {
		if errors.Is(err, storageErr.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, "invalid credentials")
			return
		}
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("internal error: %s", err))
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("successful deleting json obj with id: %d", id.Id))
}

func (h *Handler) GetJSON(c *gin.Context) {
	var id Id

	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("not correct request: %s", err))
		return
	}

	json, err := h.Storage.GetJSON(id.Id)
	if err != nil {
		if errors.Is(err, storageErr.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, "invalid credentials")
			return
		}
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("internal error: %s", err))
		return
	}

	c.JSON(http.StatusOK, map[string]string{"json": json})
}
