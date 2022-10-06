package db

import (
	"context"
	"textBoard/schema"
)

type Repository interface {
	Close()
	InsertPost(ctx context.Context, post schema.Post) error
	ListPosts(ctx context.Context, skip uint64, take uint64) ([]schema.Post, error)
}

// Инверсия управления

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func InsertPost(ctx context.Context, post schema.Post) error {
	return impl.InsertPost(ctx, post)
}

func ListPosts(ctx context.Context, skip uint64, take uint64) ([]schema.Post, error) {
	return impl.ListPosts(ctx, skip, take)
}
