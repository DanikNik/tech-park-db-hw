package controllers

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"net/http"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/models"
)

func PostGetOne(ctx *routing.Context) error { return nil }

func PostUpdate(ctx *routing.Context) error { return nil }

func PostsCreate(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	posts := []models.Post{}
	json.Unmarshal(ctx.PostBody(), &posts)
	newPosts, err := db.CreatePostsBulk(slugOrId, posts)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(http.StatusNotFound)
		} else if err == db.ErrConflict {
			ctx.SetStatusCode(http.StatusConflict)
		} else {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return nil
		}
		data, _ := json.Marshal(models.NewErrorMessage())
		ctx.Write(data)
		return nil
	}
	ctx.SetStatusCode(http.StatusCreated)
	data, _ := json.Marshal(newPosts)
	ctx.Write(data)
	return nil
}
