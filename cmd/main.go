package main

import (
	app "github.com/SergeyButGusaim/go_final_project-main"
	"github.com/SergeyButGusaim/go_final_project-main/pkg/handler"
	"github.com/SergeyButGusaim/go_final_project-main/pkg/service"
	store2 "github.com/SergeyButGusaim/go_final_project-main/pkg/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

const (
	port = "7540"
)

func main() {
	logrus.Println("Инициируем соединение с базой данных")
	dbname, err := store2.CheckingForDb()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("С базой данных %s соединились. Всё ОК", dbname)

	st := store2.NewStore(store2.GetDB(dbname))
	svc := service.NewService(st)
	h := handler.NewHandler(svc)
	srv := new(app.Server)
	err = srv.Run(port, h.InitRoutes())
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Printf("Сервер запущен на порту %s", port)

}
