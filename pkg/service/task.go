package service

import (
	"github.com/SergeyButGusaim/go_final_project/pkg/model"
	"github.com/SergeyButGusaim/go_final_project/pkg/store"
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
func (t *TaskService) GetTasks(search string) (model.ListTasks, error) {
	return t.stor.GetTasks(search)
}
func (t *TaskService) GetTaskById(id string) (model.Task, error) {
	return t.stor.GetTaskById(id)
}
func (t *TaskService) UpdateTask(task model.Task) error {
	return t.stor.UpdateTask(task)
}
func (t *TaskService) DeleteTask(id string) error {
	return t.stor.DeleteTask(id)
}
func (t *TaskService) TaskDone(id string) error {
	return t.stor.TaskDone(id)
}
