package controllers

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/models"
)

func postIDToInt(ctx *routing.Context) int {
	id, _ := strconv.Atoi(ctx.Param("id"))
	return int(id)
}

func PostGetOne(ctx *routing.Context) error {
	postFull := &models.PostFull{}
	postFull.Post = &models.Post{}

	postFull.Post.Id = int(postIDToInt(ctx))
	related := ctx.QueryArgs().Peek("related")
	err := db.GetPostFullData(strings.Split(string(related), ","), postFull)
	if err == db.ErrNotFound {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.Write(models.ErrorByteMessage)
		return nil
	}
	data, _ := json.Marshal(postFull)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(data)
	return nil
}

func PostUpdate(ctx *routing.Context) error {
	post := &models.Post{}
	post.Id = int(postIDToInt(ctx))

	pU := &models.PostUpdate{}
	json.Unmarshal(ctx.PostBody(), pU)
	err := db.UpdatePost(post, pU)
	if err == db.ErrNotFound {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.Write(models.ErrorByteMessage)
		return nil
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(post)
	ctx.Write(data)
	return nil

}

func PostsCreate(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	posts := models.Posts{}
	json.Unmarshal(ctx.PostBody(), &posts)
	newPosts, err := db.CreatePostsBulk(slugOrId, &posts)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(http.StatusNotFound)
		} else if err == db.ErrConflict {
			ctx.SetStatusCode(http.StatusConflict)
		} else {
			ctx.SetStatusCode(http.StatusInternalServerError)
			return nil
		}
		ctx.Write(models.ErrorByteMessage)
		return nil
	}
	ctx.SetStatusCode(http.StatusCreated)
	data, _ := json.Marshal(newPosts)
	ctx.Write(data)
	return nil
}
