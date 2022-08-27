package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"github.com/bgoldovsky/shortener/internal/app/models"
	internalErrors "github.com/bgoldovsky/shortener/internal/app/repositories/urls/errors"
)

const timeout = time.Second * 3

var statement = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type database interface {
	PingContext(ctx context.Context) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Close() error
	Begin() (*sql.Tx, error)
}

type postgresRepository struct {
	db database
}

func NewRepository(dsn string) (*postgresRepository, error) {
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

	return &postgresRepository{
		db: db,
	}, nil
}

// Add Сохраняет URL
func (r *postgresRepository) Add(ctx context.Context, urlID, url, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query, args, err := buildAddQuery(urlID, url, userID)
	if err != nil {
		return fmt.Errorf("build add url query error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			query, args, err = buildGetIDQuery(url)
			if err != nil {
				return fmt.Errorf("build get url id query error: %w", err)
			}

			err = r.db.QueryRowContext(ctx, query, args...).Scan(&urlID)
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

func buildAddQuery(urlID, url, userID string) (sql string, args []interface{}, err error) {
	q := statement.
		Insert("urls").
		Columns("id,url,user_id").
		Values(urlID, url, userID)

	return q.ToSql()
}

func buildGetIDQuery(url string) (sql string, args []interface{}, err error) {
	q := statement.
		Select("id").
		From("urls").
		Where(
			sq.Eq{"url": url},
			sq.Eq{"deleted_at": nil},
		)

	return q.ToSql()
}

// AddBatch Сохраняет список URL
// Можно было бы добавить данные одним запросом, но в рамках урока хочется попробовать транзакции
func (r *postgresRepository) AddBatch(ctx context.Context, urls []models.URL, userID string) error {
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
func (r *postgresRepository) Get(ctx context.Context, urlID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query, args, err := buildGetQuery(urlID)
	if err != nil {
		return "", fmt.Errorf("build get url query error: %w", err)
	}

	var (
		url       sql.NullString
		deletedAt sql.NullTime
	)

	_ = r.db.QueryRowContext(ctx, query, args...).Scan(&url, &deletedAt)
	if deletedAt.Valid {
		return "", internalErrors.ErrURLDeleted
	}
	if !url.Valid {
		return "", internalErrors.ErrURLNotFound
	}

	return url.String, nil

}

func buildGetQuery(urlID string) (sql string, args []interface{}, err error) {
	q := statement.
		Select("url", "deleted_at").
		From("urls").
		Where(sq.And{
			sq.Eq{"id": urlID},
		})

	return q.ToSql()
}

// GetList Возвращает список всех сокращенных URL
func (r *postgresRepository) GetList(ctx context.Context, userID string) ([]models.URL, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	query, args, err := buildGetListQuery(userID)
	if err != nil {
		return nil, fmt.Errorf("build get urls query error: %w", err)
	}

	res := make([]models.URL, 0)
	rows, _ := r.db.QueryContext(ctx, query, args...)
	err = rows.Err()
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

func buildGetListQuery(userID string) (sql string, args []interface{}, err error) {
	q := statement.
		Select("id, url").
		From("urls").
		Where(sq.And{
			sq.Eq{"user_id": userID},
			sq.Eq{"deleted_at": nil},
		})

	return q.ToSql()
}

// Delete Удаляет список URL указанного пользователя
func (r *postgresRepository) Delete(ctx context.Context, urlsBatch []models.UserCollection) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.PrepareContext(ctx, `update urls set deleted_at=$1 where id=$2 and user_id=$3 and deleted_at is null;`)
	if err != nil {
		return err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	now := time.Now()

	for _, collection := range urlsBatch {
		for _, urlID := range collection.URLIDs {
			if _, err = stmt.ExecContext(ctx, now, urlID, collection.UserID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// Ping Проверяет доступность базы данных
func (r *postgresRepository) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return r.db.PingContext(ctx)
}

// Close Закрывает соединение
func (r *postgresRepository) Close() error {
	return r.db.Close()
}
