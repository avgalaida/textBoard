package search

import (
	"context"
	"github.com/avgalaida/textBoard/schema"
)

type Repository interface {
	Close()
	InsertPost(ctx context.Context, post schema.Post) error
	SearchPosts(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Post, error)
}

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

func SearchPosts(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Post, error) {
	return impl.SearchPosts(ctx, query, skip, take)
}
