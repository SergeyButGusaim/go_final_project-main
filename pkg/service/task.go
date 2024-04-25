package service

import (
	"github.com/SergeyButGusaim/go_final_project-main/pkg/model"
	"github.com/SergeyButGusaim/go_final_project-main/pkg/store"
)

type TaskService struct {
	stor store.TodoTask
}

func NewTodoTaskService(stor store.TodoTask) *TaskService {
	return &TaskService{stor: stor}
}
func (t *TaskService) NextDate(nd model.NextDate) (string, error) {
	return t.stor.NextDate(nd)
}
func (t *TaskService) CreateTask(task model.Task) (int64, error) {
	return t.stor.CreateTask(task)
}

func (t *TaskService) GetTaskById(id string) (model.Task, error) {
	return t.stor.GetTaskById(id)
}
