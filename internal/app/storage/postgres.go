package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/koteyye/shortener/internal/app/models"
	"go.uber.org/zap"
)

// DBStorage структура БД.
type DBStorage struct {
	db *sqlx.DB
}

// NewPostgres возвращает новый экземпляр БД.
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
	_, err = db.ExecContext(ctx, `create table if not exists shorturl
	(
		id           serial       not null primary key,
		shortURL     varchar(256) not null,
		originalURL  varchar      not null unique,
		user_created uuid         not null,
		is_deleted    bool         not null default false
	);`)
	if err != nil {
		logger.Fatal(err.Error(), "event", "create table and index")
	}
}

// GetURLByUser получить список URL, созданных текущим пользователем.
func (d *DBStorage) GetURLByUser(ctx context.Context, userID string) ([]*models.URLList, error) {
	var result []*models.URLList
	err := d.db.SelectContext(ctx, &result, "select originalURL, shortURL from shorturl where user_created = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("не удалось получить сокращенный url из бд: %v", err)
	}
	return result, nil
}

// AddURL добавить URL в базу.
func (d *DBStorage) AddURL(ctx context.Context, shortURL string, originalURL string, userID string) error {
	_, err := d.db.ExecContext(ctx, "insert into shorturl (shorturl, originalurl, user_created) values ($1, $2, $3)", shortURL, originalURL, userID)
	if err != nil {
		return fmt.Errorf("can't add URL to DB: %w", err)
	}
	return nil
}

// GetURL получить URL из базы.
func (d *DBStorage) GetURL(ctx context.Context, shortURL string) (*models.SingleURL, error) {
	var result models.SingleURL
	err := d.db.GetContext(ctx, &result, "select shorturl, originalurl, is_deleted from shorturl where shorturl = $1", shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.SingleURL{}, models.ErrNotFound
		}
		return nil, fmt.Errorf("не удалось получить сокращенный url из бд: %w", err)
	}
	return &result, nil
}

// GetShortURL получить сокращенный URL из базы.
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

// DeleteURLByUser удалить из базы сокращенные URL по поступающему каналу.
func (d *DBStorage) DeleteURLByUser(ctx context.Context, urls chan string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при старте транзакции удаления: %v", err)
	}
	defer tx.Rollback()

	for url := range urls {
		_, err := tx.ExecContext(ctx, "update shorturl set is_deleted = true where shorturl = $1", url)
		if err != nil {
			return fmt.Errorf("ошибка при обновлении записи с shorturl: %s\n ошибка: %v", url, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetDBPing проверить подключение к БД.
func (d *DBStorage) GetDBPing(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
