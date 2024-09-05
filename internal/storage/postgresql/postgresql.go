package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"jps/internal/config"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

var (
	jsonsTable = "jsons"
)

var (
	ErrInternal           = errors.New("internal error")
	ErrInvalidCredentials = errors.New("invalid credentials")
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

	stmt, err := psql.db.Prepare(fmt.Sprintf("INSERT INTO %s (json) values ($1::json) RETURNING id", jsonsTable))
	slog.Info(json)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	result := stmt.QueryRow(json)

	if err := result.Scan(&idRes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return idRes, nil
}

func (psql *PostgreDB) GetJSON(id int) (json string, err error) {
	var jsonRes string

	stmt, err := psql.db.Prepare(fmt.Sprintf("SELECT json FROM %s WHERE id=$1", jsonsTable))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInternal, err)
	}

	result := stmt.QueryRow(id)

	if err := result.Scan(&jsonRes); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
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

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	} else if count == 0 {
		return ErrInvalidCredentials
	}

	return nil
}
