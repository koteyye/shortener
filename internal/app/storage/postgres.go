package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

const (
	PqDuplicateErr = "23505"
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

	logger := *newLogger.Sugar()
	_, err = db.ExecContext(ctx, `create table if not exists shorturl (
    id serial not null primary key ,
    shortURL varchar(256) not null,
    originalURL varchar not null unique
)`)
	if err != nil {
		logger.Fatal(err.Error(), "event", "create table and index")
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

func (d *DBStorage) AddURL(ctx context.Context, s string, s2 string) error {
	_, err := d.db.ExecContext(ctx, "insert into shorturl (shorturl, originalurl) values ($1, $2)", s, s2)
	if err != nil {
		return fmt.Errorf("can't add URL to DB: %w", err)
	}
	return nil
}

func (d *DBStorage) GetURL(ctx context.Context, s string) (string, error) {
	var result string
	if err := d.db.GetContext(ctx, &result, "select originalurl from shorturl where shorturl = $1", s); err != nil {
		return "", fmt.Errorf("can't select longURL: %w", err)
	}
	return result, nil
}

func (d *DBStorage) GetShortURL(ctx context.Context, s string) (string, error) {
	var result string
	if err := d.db.GetContext(ctx, &result, "select shorturl from shorturl where originalurl = $1", s); err != nil {
		return "", fmt.Errorf("can't select shortURL: %w", err)
	}
	return result, nil
}

func (d *DBStorage) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
