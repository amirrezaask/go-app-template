package storage

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	conn *sql.DB
}

func NewMySQL(connectionString string) (*MySQL, error) {
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return &MySQL{
		conn: conn,
	}, nil

}
