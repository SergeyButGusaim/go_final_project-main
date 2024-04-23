package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/SergeyButGusaim/go_final_project-main/pkg/model"
	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"error"`
}

func NewResponseError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, Error{Message: message})
}

func (h Handler) getNextDate(c *gin.Context) {
	var nd model.NextDate

	err := c.ShouldBindQuery(&nd)
	if err != nil {
		log.Println(err)
		c.Data(400, "text/plain", []byte(err.Error()))
		return
	}
	str, err := h.service.NextDate(nd)
	if err != nil {
		log.Println(err)
		c.Data(400, "text/plain", []byte(err.Error()))
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte(str))
}

func (h *Handler) createTask(c *gin.Context) {
	var task model.Task
	if c.ShouldBindJSON(&task) == nil {
		logrus.Println(fmt.Sprintf(
			"Получили объект Task со следующими данными: date: %s, title: %s, comment: %s, repeat: %s",
			task.Date, task.Title, task.Comment, task.Repeat))
	}
	id, err := h.service.TodoTask.CreateTask(task)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{"id": id})

}
