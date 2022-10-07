package search

import (
	"context"
	"encoding/json"
	"github.com/avgalaida/textBoard/schema"
	"github.com/olivere/elastic"
	"log"
)

type ElasticRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

func (r *ElasticRepository) Close() {
}

func (r *ElasticRepository) InsertPost(ctx context.Context, post schema.Post) error {
	_, err := r.client.Index().
		Index("posts").
		Type("post").
		Id(post.ID).
		BodyJson(post).
		Refresh("wait_for").
		Do(ctx)
	return err
}

func (r *ElasticRepository) SearchPosts(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Post, error) {
	result, err := r.client.Search().
		Index("posts").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	posts := []schema.Post{}
	for _, hit := range result.Hits.Hits {
		var post schema.Post
		if err = json.Unmarshal(*hit.Source, &post); err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}