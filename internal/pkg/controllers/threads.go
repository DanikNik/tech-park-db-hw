package controllers

import (
	"encoding/json"
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/models"
)

func ThreadCreate(ctx *routing.Context) error {
	threadData := models.Thread{}
	forumSlug := ctx.Param("slug")
	json.Unmarshal(ctx.PostBody(), &threadData)
	threadData.Forum = forumSlug
	err := db.CreateThread(&threadData)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.ErrorByteMessage)
			return nil
		} else if err == db.ErrConflict {
			ctx.SetStatusCode(fasthttp.StatusConflict)
			existingThread, err := db.GetThread(threadData.Slug)
			fmt.Printf("%+v", existingThread)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return nil
			}
			data, _ := json.Marshal(existingThread)
			ctx.Write(data)
			return nil
		}
	}
	data, _ := json.Marshal(&threadData)
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write(data)
	return nil
}

func ThreadGetOne(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	threadData, err := db.GetThread(slugOrId)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.ErrorByteMessage)
			return nil
		} else {
			ctx.SetStatusCode(500)
			ctx.Write([]byte(err.Error()))
			return nil
		}
	}
	data, _ := json.Marshal(threadData)
	ctx.Write(data)
	return nil
}

func ThreadGetPosts(ctx *routing.Context) error {

	posts := models.Posts{}

	err := db.GetPostsByThread(ctx.Param("slug_or_id"),
		ctx.QueryArgs().GetUintOrZero("limit"), getBooleanFromQueryParam("desc", ctx.QueryArgs()),
		ctx.QueryArgs().GetUintOrZero("since"),
		string(ctx.QueryArgs().Peek("sort")), &posts)

	if err == db.ErrNotFound {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write(models.ErrorByteMessage)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(&posts)
	ctx.Write(data)

	return nil
}

func ThreadUpdate(ctx *routing.Context) error {
	threadSlugOrId := ctx.Param("slug_or_id")
	threadUpdateData := models.ThreadUpdate{}
	json.Unmarshal(ctx.PostBody(), &threadUpdateData)
	threadData, err := db.UpdateThread(threadSlugOrId, &threadUpdateData)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.ErrorByteMessage)
			return nil
		}
	}
	data, _ := json.Marshal(threadData)
	ctx.Write(data)
	return nil
}

func ThreadVote(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	voteData := models.Vote{}
	json.Unmarshal(ctx.PostBody(), &voteData)
	threadData, err := db.DoVote(slugOrId, voteData)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.Write(models.ErrorByteMessage)
			return nil
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil
	}
	data, _ := json.Marshal(threadData)
	ctx.Write(data)
	return nil
}
