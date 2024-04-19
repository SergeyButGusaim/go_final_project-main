package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h Handler) getNextDate(c *gin.Context) {
	nowStr := c.Param("now")
	date := c.Param("date")
	repeat := c.Param("repeat")

	now, err := time.Parse(DATE_FORMAT, nowStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты"})
		return
	}

	next, err := h.service.NextDate(now, date, repeat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"next": next})
}
