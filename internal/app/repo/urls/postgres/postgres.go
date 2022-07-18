package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repo/urls/errors"
)

const timeout = time.Second * 3

type database interface {
	PingContext(ctx context.Context) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Close() error
	Begin() (*sql.Tx, error)
}

type pgRepo struct {
	db database
}

func NewRepo(dsn string) (*pgRepo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)

	query := `create table if not exists urls 
(
    id varchar(10) not null primary key,
    url varchar(500) not null unique,
    user_id varchar(10) not null,
    created_at timestamp with time zone default now() not null,
    deleted_at  timestamp with time zone default null
);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &pgRepo{
		db: db,
	}, nil
}

// Add Сохраняет URL
func (r *pgRepo) Add(ctx context.Context, urlID, url, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `insert into urls(id,url,user_id) values ($1,$2,$3)`, urlID, url, &userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			err = r.db.QueryRowContext(ctx, "select id from urls where url=$1", url).Scan(&urlID)
			if err != nil {
				return err
			}
			err = internalErrors.NewNotUniqueURLErr(urlID, url, err)
			return err
		}

		return err
	}

	return nil
}

// AddBatch Сохраняет список URL
// Можно было бы добавить данные одним запросом, но в рамках урока хочется попробовать транзакции
func (r *pgRepo) AddBatch(ctx context.Context, urls []models.URL, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.PrepareContext(ctx, `insert into urls(id,url,user_id) values ($1,$2,$3);`)
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	for idx := range urls {
		if _, err = stmt.ExecContext(ctx, urls[idx].ShortURL, urls[idx].OriginalURL, userID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Get Возвращает URL
func (r *pgRepo) Get(ctx context.Context, urlID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var url sql.NullString
	_ = r.db.QueryRowContext(ctx, `select url from urls where id=$1 and deleted_at is null`, urlID).Scan(&url)
	if url.Valid {
		return url.String, nil
	}

	return "", internalErrors.ErrURLNotFound
}

// GetList Возвращает список всех сокращенных URL
func (r *pgRepo) GetList(ctx context.Context, userID string) ([]models.URL, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res := make([]models.URL, 0)
	rows, _ := r.db.QueryContext(ctx, `select id, url from urls where user_id=$1 and deleted_at is null;`, userID)
	err := rows.Err()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var url models.URL
		err = rows.Scan(&url.ShortURL, &url.OriginalURL)
		if err != nil {
			return nil, err
		}

		res = append(res, url)
	}

	return res, nil
}

// Ping Проверяет доступность базы данных
func (r *pgRepo) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return r.db.PingContext(ctx)
}

func (r *pgRepo) Close() error {
	return r.db.Close()
}
