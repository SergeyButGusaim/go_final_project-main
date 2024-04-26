package service

import (
	"github.com/SergeyButGusaim/go_final_project-main/pkg/model"
	"github.com/SergeyButGusaim/go_final_project-main/pkg/store"
)

type TodoTask interface {
	NextDate(nd model.NextDate) (string, error)
	CreateTask(task model.Task) (int64, error)
	GetTaskById(id string) (model.Task, error)
	GetTasks(search string) (model.ListTasks, error)
	UpdateTask(task model.Task) error
	TaskDone(id string) error
	DeleteTask(id string) error
}

type Service struct {
	TodoTask
}

func NewService(store *store.Store) *Service {
	return &Service{
		TodoTask: NewTodoTaskService(store.TodoTask),
	}
}
