package storage

import (
	"errors"
	"fmt"
	"jps/internal/config"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

var (
	jsonsTable = "jsons"
)

var (
	ErrInternal = errors.New("Internal error")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)
type PostgreDB struct {
	db *sqlx.DB
}

func NewPostgreDB(db *sqlx.DB) *PostgreDB {
	return &PostgreDB{db: db}
}

func NewSqlxDB(cfg config.Config) (*sqlx.DB, error) {
	op := "storage.NewSqlxDB"
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.PortDB, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLmode))
	if err != nil {
		return nil, fmt.Errorf("%w:%s", err, op)
	}
	// "host=localhost port=5432 user=user password=password dbname=jpsdb sslmode=disable"
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%w:%s", err, op)
	}
	return db, nil
}

func (psql *PostgreDB) NewJSON(json string) (id int, err error) {
	var idRes int

	stmt, err := psql.db.Prepare(fmt.Sprintf("INSERT INTO %s (data) values ($2::json) RETURNING id", jsonsTable))
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	result := stmt.QueryRow(json)

	if err := result.Scan(&idRes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
		}
		return 0, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return idRes, nil
}

func (psql *PostgreDB) GetJSON(id int) (json string, err error) {
	var jsonRes string

	stmt, err := psql.db.Prepare(fmt.Sprintf("SELECT json FROM TABLE %s WHERE id=$2", jsonsTable))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	result := stmt.QueryRow(jsonsTable, id)

	if err := result.Scan(&jsonRes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
		}
		return "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return jsonRes, nil

}

func (psql *PostgreDB) DeleteJSON(id int) (err error) {
	stmt, err := psql.db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id=$1", jsonsTable))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %w", ErrInvalidCredentials, err)
		}
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return nil
}
