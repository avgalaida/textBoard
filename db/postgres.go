package db

import (
	"context"
	"database/sql"
	"github.com/avgalaida/textBoard/schema"
	_ "github.com/lib/pq"
	"log"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		db,
	}, nil
}

func (r *PostgresRepository) Close() {
	if err := r.db.Close(); err != nil {
		log.Fatal(err)
	}
}
func (r *PostgresRepository) InsertPost(ctx context.Context, post schema.Post) error {
	_, err := r.db.Exec("INSERT INTO posts(id, body, created_at) VALUES($1, $2, $3)", post.ID, post.Body, post.CreatedAt)
	return err
}

func (r *PostgresRepository) ListPosts(ctx context.Context, skip uint64, take uint64) ([]schema.Post, error) {
	rows, err := r.db.Query("SELECT * FROM posts ORDER BY id DESC OFFSET $1 LIMIT $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var posts []schema.Post
	for rows.Next() {
		post := schema.Post{}
		if err = rows.Scan(&post.ID, &post.Body, &post.CreatedAt); err == nil {
			posts = append(posts, post)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
