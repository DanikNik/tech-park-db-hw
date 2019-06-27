package controllers

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/models"
)

func getBooleanFromQueryParam(k string, args *fasthttp.Args) bool {
	v := args.Peek(k)
	if v != nil && v[0] == 't' {
		return true
	}
	return false
}

func ForumCreate(ctx *routing.Context) error {
	var forumData models.Forum
	json.Unmarshal(ctx.PostBody(), &forumData)
	err := db.CreateForum(&forumData)
	switch err {
	case db.ErrNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.ErrorByteMessage)
		return nil
	case db.ErrConflict:
		existingForum, err := db.GetForum(forumData.Slug)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return nil
		}
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(existingForum)
		ctx.Write(data)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	data, _ := json.Marshal(forumData)
	ctx.Write(data)
	return nil
}

func ForumGetOne(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	forumData, err := db.GetForum(slug)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.ErrorByteMessage)
			return nil
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return nil
		}
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(forumData)
	ctx.Write(data)
	return nil
}

func ForumGetThreads(ctx *routing.Context) error {
	threads, err := db.GetThreadsByForum(
		ctx.Param("slug"),
		ctx.QueryArgs().GetUintOrZero("limit"),
		getBooleanFromQueryParam("desc", ctx.QueryArgs()),
		string(ctx.QueryArgs().Peek("since")),
	)

	if err == db.ErrNotFound {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.ErrorByteMessage)
		return nil
	}
	data, _ := json.Marshal(threads)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write(data)
	return nil
}

func ForumGetUsers(ctx *routing.Context) error {
	userList, err := db.GetUsersByForum(
		ctx.Param("slug"),
		ctx.QueryArgs().GetUintOrZero("limit"),
		getBooleanFromQueryParam("desc", ctx.QueryArgs()),
		string(ctx.QueryArgs().Peek("since")),
	)

	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.ErrorByteMessage)
			return nil
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(&userList)
	ctx.Write(data)
	return nil
}
