package storage

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

type DBStorage struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (d *DBStorage) AddURL(s string, s2 string) error {
	//TODO implement me
	return errors.New("not implement")
}

func (d *DBStorage) GetURL(s string) (string, error) {
	//TODO implement me
	return "", errors.New("not implement")
}

func (d *DBStorage) Ping() error {
	return d.db.Ping()
}
