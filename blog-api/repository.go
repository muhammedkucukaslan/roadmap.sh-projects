package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
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
	query := `CREATE TABLE IF NOT EXISTS blog (
					id SERIAL,
					title VARCHAR(100),
					content VARCHAR,
					category VARCHAR(100),
					tags 	TEXT[],
					created_at TIMESTAMP DEFAULT NOW(), 
					updated_at TIMESTAMP DEFAULT NOW()
	 		  );`

	_, err := r.db.Exec(query)
	return err
}
func (r *Repository) Create(data CreateBlogRequest) error {

	query := `INSERT INTO blog (title, content, category, tags) values 
	($1, $2, $3, $4)`

	_, err := r.db.Exec(query,
		data.Title,
		data.Content,
		data.Category,
		pq.Array(data.Tags),
	)
	return err
}
func (r *Repository) Update(id int, req UpdateBlogRequest) error {

	blog, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if req.Title != nil {
		blog.Title = *req.Title
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}

	if req.Category != nil {
		blog.Category = *req.Category
	}

	if req.Tags != nil {
		blog.Tags = *req.Tags
	}

	query := `UPDATE blog
        SET title = $2, content = $3, category = $4, tags = $5, updated_at = $6
        WHERE id = $1`

	_, err = r.db.Exec(query, blog.ID, blog.Title, blog.Content, blog.Category, pq.Array(blog.Tags), time.Now())
	if err != nil {
		return fmt.Errorf("failed to update blog with ID %d: %w", id, err)
	}

	return nil
}
func (r *Repository) Delete(id int) error {
	query := `DELETE FROM blog WHERE   id = $1`

	resp, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}
	affected, _ := resp.RowsAffected()
	if affected == 0 {
		return errInvalidID
	}
	return nil
}
func (r *Repository) GetByID(id int) (Blog, error) {

	query := `SELECT * FROM blog WHERE id = $1`

	var blog Blog
	err := r.db.QueryRow(query, id).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.Category,
		pq.Array(&blog.Tags),
		&blog.CreatedAt,
		&blog.UpdatedAt)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return Blog{}, errInvalidID
		}

		return Blog{}, err

	}

	return blog, nil

}

func (r *Repository) GetBlogs() ([]Blog, error) {

	query := `SELECT * FROM blog`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog
	if rows.Next() {
		var blog Blog

		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.Category,
			pq.Array(&blog.Tags),
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}
	fmt.Println(blogs)
	return blogs, nil
}
