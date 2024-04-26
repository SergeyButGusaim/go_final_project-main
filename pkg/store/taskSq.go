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
	"github.com/sirupsen/logrus"
)

const (
	ErrInvalidDateFormat        = "Неверный формат даты"
	ErrEmptyRepeatRule          = "Отсутствует правило повторения"
	ErrInvalidRepeatRuleD       = "Неверный формат правила повторения 'd'"
	ErrInvalidNumberOfDays      = "Неверный формат количества дней"
	ErrInvalidNumberOfDaysRange = "Неверное количество дней"
	ErrUnsupportedRepeatRule    = "Неподдерживаемое правило повторения"
	dbname                      = "scheduler"
	MaxLimit                    = 25
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

	query := fmt.Sprintf("INSERT INTO %s (title, comment, date, repeat) VALUES ($1, $2, $3, $4)", dbname)
	res, err := t.db.Exec(query, task.Title, task.Comment, task.Date, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
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

func (t *TaskSq) GetTasks(search string) (model.ListTasks, error) {
	var tasks []model.Task
	var query string

	switch typeSearch(search) {
	case 0:
		query = fmt.Sprintf("SELECT * FROM %s ORDER BY date LIMIT ?", dbname)
		err := t.db.Select(&tasks, query, MaxLimit)
		if err != nil {
			return model.ListTasks{}, err
		}
	case 1:
		s, _ := time.Parse(`02.01.2006`, search)
		st := s.Format(`20060102`)
		query = fmt.Sprintf("SELECT * FROM %s WHERE date = ? ORDER BY date LIMIT ?", dbname)
		err := t.db.Select(&tasks, query, st, MaxLimit)
		if err != nil {
			return model.ListTasks{}, err
		}
	case 2:
		searchQuery := fmt.Sprintf("%%%s%%", search)
		query := `SELECT * FROM scheduler WHERE LOWER(title) LIKE ? OR LOWER(comment) LIKE ? ORDER BY date LIMIT ?`
		rows, err := t.db.Queryx(query, searchQuery, searchQuery, MaxLimit)
		if err != nil {
			return model.ListTasks{}, err
		}
		for rows.Next() {
			var task model.Task
			err := rows.StructScan(&task)
			if err != nil {
				return model.ListTasks{}, err
			}
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		return model.ListTasks{Tasks: []model.Task{}}, nil
	}
	return model.ListTasks{Tasks: tasks}, nil
}

func typeSearch(str string) int {
	if str == "" {
		return 0
	}
	_, err := time.Parse(`02.01.2006`, str)
	if err == nil {
		return 1
	}
	return 2
}

func (t *TaskSq) UpdateTask(task model.Task) error {
	err := t.checkTask(&task)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET title = ?, comment = ?, date = ?, repeat = ? WHERE id = ?", dbname)
	_, err = t.db.Exec(query, task.Title, task.Comment, task.Date, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}
func (t *TaskSq) DeleteTask(id string) error {
	_, err := t.GetTaskById(id)
	if err != nil {
		return err
	}
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE id = ?", dbname)
	_, err = t.db.Exec(queryDelete, id)
	if err != nil {
		return err
	}
	return nil
}
func (t *TaskSq) TaskDone(id string) error {
	task, err := t.GetTaskById(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		queryDeleteTask := fmt.Sprintf("DELETE FROM %s WHERE id = ?", dbname)
		logrus.Println(queryDeleteTask)
		t.db.Exec(queryDeleteTask, id)
		return nil
	}

	nd := model.NextDate{
		Date:   task.Date,
		Now:    time.Now().Format(`20060102`),
		Repeat: task.Repeat,
	}

	newDate, err := t.NextDate(nd)
	if err != nil {
		return err
	}

	task.Date = newDate
	queryUpdateTask := fmt.Sprintf("UPDATE %s SET date = ? WHERE id = ?", dbname)
	logrus.Println(queryUpdateTask)
	_, err = t.db.Exec(queryUpdateTask, task.Date, id)
	if err != nil {
		return err
	}
	return nil

}
