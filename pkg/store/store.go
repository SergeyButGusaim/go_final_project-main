package store

import (
	"github.com/SergeyButGusaim/go_final_project-main/pkg/model"
	"github.com/jmoiron/sqlx"
)

type TodoTask interface {
	NextDate(nd model.NextDate) (string, error)
	CreateTask(task model.Task) (int64, error)
	GetTaskById(id string) (model.Task, error)
	GetTasks(search string) (model.ListTasks, error)
	UpdateTask(task model.Task) error
	DeleteTask(id string) error
	TaskDone(id string) error
}

type Store struct {
	TodoTask
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		TodoTask: NewTaskSq(db),
	}
}
