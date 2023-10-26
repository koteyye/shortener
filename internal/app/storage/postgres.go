package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/shortener/internal/app/models"
	"go.uber.org/zap"
	"time"
)

type ctxUserKEy string

const userIDKey = ctxUserKEy("userId")

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
    originalURL varchar not null unique,
    user_created uuid not null
)`)
	if err != nil {
		logger.Fatal(err.Error(), "event", "create table and index")
	}
}

func (d *DBStorage) GetURLByUser(ctx context.Context, userId string) ([]*models.AllURLs, error) {
	var result []*models.AllURLs
	err := d.db.SelectContext(ctx, &result, "select originalURL, shortURL from shorturl where user_created = $1", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("не удалось получить сокращенный url из бд: %v", err)
	}
	return result, nil
}

func (d *DBStorage) AddURL(ctx context.Context, shortURL string, originalURL string) error {
	userID := ctx.Value(userIDKey)
	_, err := d.db.ExecContext(ctx, "insert into shorturl (shorturl, originalurl, user_created) values ($1, $2, $3)", shortURL, originalURL, userID)
	if err != nil {
		return fmt.Errorf("can't add URL to DB: %w", err)
	}
	return nil
}

func (d *DBStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	var result string
	err := d.db.GetContext(ctx, &result, "select originalurl from shorturl where shorturl = $1", shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNotFound
		}
		return "", fmt.Errorf("не удалось получить сокращенный url из бд: %w", err)
	}
	return result, nil
}

func (d *DBStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	var result string
	err := d.db.GetContext(ctx, &result, "select shorturl from shorturl where originalurl = $1", originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrNotFound
		}
		return "", fmt.Errorf("не удалось получить сокращенный url из бд: %w", err)
	}
	return result, nil
}

func (d *DBStorage) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
