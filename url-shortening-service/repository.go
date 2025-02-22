package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connStr string) (*Repository, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil

}

func (r *Repository) Init() error {
	query := `CREATE TABLE IF NOT EXISTS urls (
	id SERIAL PRIMARY KEY,
	url VARCHAR(250) NOT NULL,		
	short_code VARCHAR(10) NOT NULL UNIQUE,
	access_count INTEGER DEFAULT 0,
	created_at TIMESTAMP WITH  TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH  TIME ZONE DEFAULT CURRENT_TIMESTAMP
)`
	_, err := r.db.Exec(query)
	return err
}

func (r *Repository) Save(url, code string) error {

	queryCreation := `INSERT INTO urls (url, short_code) VALUES ($1,$2)`
	_, err := r.db.Exec(queryCreation, url, code)

	return err
}

func (r *Repository) Find(code string) (UrlObject, error) {
	query := `SELECT * FROM urls WHERE short_code = $1`

	var urlObject UrlObject

	row := r.db.QueryRow(query, code)

	err := row.Scan(&urlObject.ID, &urlObject.Url, &urlObject.ShortCode, &urlObject.AccessCount, &urlObject.CreatedAt, &urlObject.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return UrlObject{}, errUrlNotFound
		}
		return UrlObject{}, err
	}
	return urlObject, nil
}

func (r *Repository) IncreaseAccessCount(code string) error {
	query := `UPDATE urls SET access_count = access_count + 1 WHERE short_code = $1`
	_, err := r.db.Exec(query, code)

	return err
}

func (r *Repository) Update(url, code string) error {
	query := `UPDATE urls SET url= $1 ,access_count = 0, updated_at = CURRENT_TIMESTAMP WHERE short_code = $2`

	_, err := r.db.Exec(query, url, code)

	return err
}

func (r *Repository) Delete(code string) error {
	query := `DELETE FROM urls WHERE code = $1`

	result, err := r.db.Exec(query, code)
	if err != nil {
		return errServerError
	}
	num, _ := result.RowsAffected()
	if num == 0 {
		return errInvalidCode
	}

	return nil
}
