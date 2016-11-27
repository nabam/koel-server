package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nabam/koel-server/models"
	"github.com/satori/go.uuid"
)

type sqliteStorage struct {
	db *sql.DB
}

func (storage sqliteStorage) GetServices() ([]models.Service, error) {
	svcs := make([]models.Service, 0)

	rows, err := storage.db.Query("SELECT service_id, service_name FROM services")
	if err != nil {
		if err.Error() == "EOF" {
			return svcs, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i string
		var n string
		err = rows.Scan(&i, &n)
		if err != nil {
			return svcs, err
		}
		id, err := uuid.FromString(i)
		if err != nil {
			return svcs, err
		}
		svcs = append(svcs, models.Service{Id: id, Name: n})
	}
	return svcs, nil
}

func (storage sqliteStorage) CreateService(name string) error {
	svc := models.NewService(name)
	_, err := storage.db.Exec("INSERT INTO services(service_id, service_name) VALUES(?, ?)", svc.Id.String(), svc.Name)
	return err
}

func (storage sqliteStorage) DeleteService(name string) error {
	_, err := storage.db.Exec("DELETE FROM services WHERE service_name = ?", name)
	return err
}

func (storage sqliteStorage) Close() error {
	return storage.db.Close()
}

func InitSqlite(fileName string) (Storage, error) {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	initQuery := "CREATE TABLE IF NOT EXISTS services(service_id TEXT PRIMARY KEY, service_name TEXT NOT NULL);"
	initQuery += "CREATE UNIQUE INDEX IF NOT EXISTS services_name ON services(service_name);"

	_, err = db.Exec(initQuery)
	if err != nil {
		return nil, err
	}

	var storage sqliteStorage
	storage.db = db
	return &storage, nil
}
