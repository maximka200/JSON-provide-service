package storage

import (
	"fmt"
	"jps/internal/config"

	"github.com/jmoiron/sqlx"
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
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.PortDB, cfg.Username, cfg.DBName, cfg.SSLmode))
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
	panic("dont impl")
}

func (psql *PostgreDB) GetJSON(id int) (json string, err error) {
	panic("dont impl")
}

func (psql *PostgreDB) DeleteJSON(id int) (err error) {
	panic("dont impl")
}
