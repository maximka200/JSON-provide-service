package handler

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Id struct {
	id int
}

// слой десериализации данных из http req

func (h *Handler) NewJSON(c *gin.Context) {
	json, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "cannot read message body")
		return
	}
	defer c.Request.Body.Close()

	jsonStr := string(json)

	id, err := h.Storage.NewJSON(jsonStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "cannot get json id")
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

	if err := h.Storage.DeleteJSON(id.id); err != nil {
		// todo: errors validation
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("internal error: %s", err))
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("successful deleting json obj with id: %d", id.id))
}

func (h *Handler) GetJSON(c *gin.Context) {
	var id Id

	if err := c.ShouldBindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("not correct request: %s", err))
		return
	}

	json, err := h.Storage.GetJSON(id.id)
	if err != nil {
		// todo: errors validation
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("internal error: %s", err))
		return
	}

	c.JSON(http.StatusOK, map[string]string{"json": json})
}
