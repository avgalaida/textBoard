package main

import (
	"github.com/avgalaida/textBoard/db"
	"github.com/avgalaida/textBoard/event"
	"github.com/avgalaida/textBoard/schema"
	"github.com/avgalaida/textBoard/util"
	"github.com/segmentio/ksuid"
	"html/template"
	"log"
	"net/http"
	"time"
)

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	post := schema.Post{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}
	if err := db.InsertPost(ctx, post); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	if err := event.PublishPostCreated(post); err != nil {
		log.Println(err)
	}

	util.ResponseOk(w, response{ID: post.ID})
}
