package store

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SergeyButGusaim/go_final_project-main/pkg/model"
	"github.com/jmoiron/sqlx"
)

const (
	ErrInvalidDateFormat        = "Неверный формат даты"
	ErrEmptyRepeatRule          = "Отсутствует правило повторения"
	ErrInvalidRepeatRuleD       = "Неверный формат правила повторения 'd'"
	ErrInvalidNumberOfDays      = "Неверный формат количества дней"
	ErrInvalidNumberOfDaysRange = "Неверное количество дней"
	ErrUnsupportedRepeatRule    = "Неподдерживаемое правило повторения"
	dbname                      = "scheduler"
)

type TaskSq struct {
	db *sqlx.DB
}

func NewTaskSq(db *sqlx.DB) *TaskSq {
	return &TaskSq{db: db}
}

func (t *TaskSq) NextDate(nd model.NextDate) (string, error) {
	nowDate, err := time.Parse("20060102", nd.Now)
	if err != nil {
		return "", errors.New(ErrInvalidDateFormat)
	}

	startDate, err := time.Parse("20060102", nd.Date)
	if err != nil {
		return "", errors.New(ErrInvalidDateFormat)
	}

	ruleParts := strings.Fields(nd.Repeat)

	if len(ruleParts) == 0 && nd.Repeat != "y" {
		return "", errors.New(ErrEmptyRepeatRule)
	}

	switch ruleParts[0] {
	case "d":
		if len(ruleParts) != 2 {
			return "", errors.New(ErrInvalidRepeatRuleD)
		}
		numOfDays, err := strconv.Atoi(ruleParts[1])
		if err != nil {
			return "", errors.New(ErrInvalidNumberOfDays)
		}
		if numOfDays <= 0 || numOfDays > 365 {
			return "", errors.New(ErrInvalidNumberOfDaysRange)
		}
		nextDate := startDate.AddDate(0, 0, numOfDays)
		for nextDate.Before(nowDate) {
			nextDate = nextDate.AddDate(0, 0, numOfDays)
		}
		return nextDate.Format("20060102"), nil
	case "y":
		nextYearDate := startDate.AddDate(1, 0, 0)
		for nextYearDate.Before(nowDate) {
			nextYearDate = nextYearDate.AddDate(1, 0, 0)
		}
		return nextYearDate.Format("20060102"), nil
	default:
		return "", errors.New(ErrUnsupportedRepeatRule)
	}
}

func (t *TaskSq) CreateTask(task model.Task) (int64, error) {
	err := t.checkTask(&task)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("INSERT INTO %s (title, comment, date, repeat) VALUES ($1, $2, $3, $4) RETURNING id", dbname)
	row := t.db.QueryRow(query, task.Title, task.Comment, task.Date, task.Repeat)

	var id int64
	if err = row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (t *TaskSq) checkTask(task *model.Task) error {
	if task.Title == "" {
		return fmt.Errorf("Ошибка. Пустое название.")
	}

	if !regexp.MustCompile(`^([wdm]\s.*|y)?$`).MatchString(task.Repeat) {
		return fmt.Errorf(ErrInvalidRepeatRuleD, "%v", task.Repeat)
	}

	now := time.Now().Format(`20060102`)

	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse(`20060102`, task.Date)
	if err != nil {
		return fmt.Errorf(ErrInvalidDateFormat)
	}

	if task.Date < now {
		if task.Repeat == "" {
			task.Date = now
		}
		if task.Repeat != "" {
			nd := model.NextDate{
				Date:   task.Date,
				Now:    now,
				Repeat: task.Repeat,
			}
			task.Date, err = t.NextDate(nd)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TaskSq) GetTaskById(id string) (model.Task, error) {
	if id == "" {
		return model.Task{}, fmt.Errorf("Не указан идентификатор")
	}
	if _, err := strconv.Atoi(id); err != nil {
		return model.Task{}, fmt.Errorf("Некорректный идентификатор")
	}
	var task model.Task
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", dbname)
	err := t.db.Get(&task, query, id)
	if err != nil {
		return model.Task{}, fmt.Errorf("Задача не найдена")
	}
	return task, err
}
