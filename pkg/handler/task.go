package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/SergeyButGusaim/go_final_project-main/pkg/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getNextDate(c *gin.Context) {
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
			"Task: date: %s, title: %s, comment: %s, repeat: %s",
			task.Date, task.Title, task.Comment, task.Repeat))
	}
	id, err := h.service.CreateTask(task)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{"id": id})

}

func (h *Handler) getTaskById(c *gin.Context) {
	id := c.Query("id")
	logrus.Println("Получен запрос на задачу с id: " + id)
	task, err := h.service.GetTaskById(id)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, task)
}

func (h *Handler) getTasks(c *gin.Context) {
	search := c.Query("search")
	logrus.Println("Получен запрос на задачи с поисковым запросом: " + search)
	list, err := h.service.TodoTask.GetTasks(search)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, list)
}

func (h *Handler) updateTask(c *gin.Context) {
	var task model.Task

	if c.ShouldBindJSON(&task) == nil {
		logrus.Println(fmt.Sprintf(
			"Получили на обновление объект task со следующими данными: "+
				"id: %s, date: %s, title: %s, comment: %s, repeat: %s",
			task.ID, task.Date, task.Title, task.Comment, task.Repeat))
	}
	_, err := h.service.TodoTask.GetTaskById(task.ID)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.TodoTask.UpdateTask(task)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{})
}
func (h *Handler) deleteTask(c *gin.Context) {
	id, _ := c.GetQuery("id")
	logrus.Println("Получен запрос на удаление задачи с id: " + id)
	err := h.service.TodoTask.DeleteTask(id)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{})
}
func (h *Handler) taskDone(c *gin.Context) {
	id, _ := c.GetQuery("id")
	logrus.Println("Получен запрос на завершение задачи с id: " + id)
	err := h.service.TodoTask.TaskDone(id)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{})
}
