package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type DBStorage struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *DBStorage {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	initDBTableShortURL(ctx, db)
	return &DBStorage{db: db}
}

func initDBTableShortURL(ctx context.Context, db *sqlx.DB) {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer newLogger.Sync()

	var tableName string
	logger := *newLogger.Sugar()
	err = db.QueryRowContext(ctx, "select table_name from information_schema.tables where table_name='shorturl';").Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		_, err := db.ExecContext(ctx, `create table shorturl (
		id serial not null primary key ,
		shortURL varchar(256) not null,
		originalURL varchar not null
	);`)
		if err != nil {
			logger.Fatalw(err.Error(), "event", "create table")
		}
		_, err = db.ExecContext(ctx, `create unique index shortURL1 on shorturl (shortURL);`)
		if err != nil {
			logger.Fatalw(err.Error(), "event", "crate index")
			return
		}
	} else {
		logger.Fatalw(err.Error())
		return
	}
}

func (d *DBStorage) CreateTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := d.db.ExecContext(ctx, `create table testshorturl 
(id serial not null primary key ,
 shortURL varchar(256) not null,
  originalURL varchar not null);`)
	if err != nil {
		return fmt.Errorf("can't build sql request create table: %v", err)
	}
	return nil
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
