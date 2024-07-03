package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error)
	Tx(ctx context.Context) (Tx, error)
}

type Tx interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, stmt string, args ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}

func Query[T FromRows](ctx context.Context, db DB, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var output []T
	for rows.Next() {
		var record T
		err := record.FromRows(rows)
		if err != nil {
			return output, err
		}

		output = append(output, record)
	}

	return output, nil
}

type FromRows interface {
	FromRows(rows *sql.Rows) error
}

func First[T FromRows](records []T, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}
	if len(records) < 1 {
		return *new(T), fmt.Errorf("record set is empty")
	}
	return records[0], nil
}
