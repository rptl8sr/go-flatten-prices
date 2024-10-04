package store

import (
	"database/sql"
	_ "embed"
	"time"

	_ "modernc.org/sqlite"
)

type store struct {
	db *sql.DB
}

type Store interface {
	GetPriceByID(id int, date string) (int, error)
}

//go:embed queries/createTable.sql
var createTable string

func initTables(db *sql.DB) error {
	_, err := db.Exec(createTable)
	return err
}

func New(path string) (Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err = initTables(db); err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}

func (s *store) GetPriceByID(id int, date string) (int, error) {
	var price int

	dateObj, err := time.Parse("060102", date)
	if err != nil {
		return 0, err
	}

	dateDB := dateObj.Format("2006-01-02")

	e := s.db.QueryRow("select price from prices where (id = ?) and (from_date = ?)", id, dateDB).Scan(&price)
	if e != nil {
		return 0, e
	}

	return price, nil
}
