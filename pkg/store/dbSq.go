package store

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func GetDB(dbname string) *sqlx.DB {
	dbname, err := CheckingForDb()
	if err != nil {
		logrus.Fatal(err)
	}
	return sqlx.MustConnect("sqlite3", dbname)
}

func CheckingForDb() (string, error) {
	dbName := "scheduler.db"

	appPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), dbName)
	logrus.Printf("Путь до файла базы данных %s", dbFile)

	_, err = os.Stat(dbFile)
	if err != nil {
		logrus.Printf("Файл базы данных %s не существует. Создаем его", dbFile)
		ok, err := dbInstallation(dbName)
		if err != nil || !ok {
			return dbName, err
		}
	}

	return dbName, nil

}

func dbInstallation(dbName string) (bool, error) {
	logrus.Println("Инициируем соединение с базой данных")
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return false, err
	}
	defer db.Close()

	createTab := `CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT, title TEXT, comment TEXT, repeat VARCHAR(128));`
	logrus.Println("Создаем таблицу scheduler")
	_, err = db.Exec(createTab)
	if err != nil {
		return false, err
	}

	createInd := `CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`
	logrus.Println("Создаем индекс")
	_, err = db.Exec(createInd)
	if err != nil {
		return false, err
	}

	logrus.Println("База данных создана")
	return true, nil
}
