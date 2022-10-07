package main

import (
	"context"
	"github.com/avgalaida/textBoard/db"
	"github.com/avgalaida/textBoard/event"
	"github.com/avgalaida/textBoard/schema"
	"github.com/avgalaida/textBoard/search"
	"github.com/avgalaida/textBoard/util"
	"log"
	"net/http"
	"strconv"
)

func onPostCreated(m event.PostCreatedMessage) {
	post := schema.Post{
		ID:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := search.InsertPost(context.Background(), post); err != nil {
		log.Println(err)
	}
}

func searchPostsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	posts, err := search.SearchPosts(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, []schema.Post{})
		return
	}

	util.ResponseOk(w, posts)
}

func listPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	posts, err := db.ListPosts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch posts")
		return
	}

	util.ResponseOk(w, posts)
}
