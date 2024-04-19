package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/SergeyButGusaim/go_final_project-main/pkg/store"
	"github.com/sirupsen/logrus"
)

const DATE_FORMAT = "20060102"

type Service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return Service{store: store}
}

func (s *Service) NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		logrus.Println("Правила повтора не указаны.")
		return "", error(nil) //TODO какую ошибку надо возвращать? Нужно ли? =)
	}
	startDate, err := time.Parse(DATE_FORMAT, date)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	param := parts[0]
	switch param {
	case "y":
		currDate := startDate.AddDate(1, 0, 0)
		for now.After(currDate) || now.Equal(currDate) {
			currDate = currDate.AddDate(1, 0, 0)
		}

		return currDate.Format(DATE_FORMAT), nil

	case "d":
		if len(parts) == 1 {
			logrus.Println("Не указан интервал в днях.")
			return "", err
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if days > 400 {
			logrus.Println("Превышен максимально допустимый интервал.")
			return "", err
		}

		currDate := startDate.AddDate(0, 0, days)
		for now.After(currDate) {
			currDate = currDate.AddDate(0, 0, days)
		}

		return currDate.Format(DATE_FORMAT), nil

	case "w":
		logrus.Println("Неподдерживаемый формат.")
		return "", nil //TODO какую ошибку надо возвращать? Нужно ли? =)
	case "m":
		logrus.Println("Неподдерживаемый формат.")
		return "", nil //TODO какую ошибку надо возвращать? Нужно ли? =)

	}
	return date, err
}
