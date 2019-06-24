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
	threadData.Forum = forumSlug
	json.Unmarshal(ctx.PostBody(), &threadData)
	err := db.CreateThread(&threadData)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			data, _ := json.Marshal(models.NewErrorMessage())
			ctx.Write(data)
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
			data, _ := json.Marshal(models.NewErrorMessage())
			ctx.Write(data)
			return nil
		}
	}
	data, _ := json.Marshal(threadData)
	ctx.Write(data)
	return nil
}

func ThreadGetPosts(ctx *routing.Context) error { return nil }

func ThreadUpdate(ctx *routing.Context) error {
	threadSlugOrId := ctx.Param("slug_or_id")
	threadUpdateData := models.ThreadUpdate{}
	json.Unmarshal(ctx.PostBody(), &threadUpdateData)
	threadData, err := db.UpdateThread(threadSlugOrId, &threadUpdateData)
	if err != nil {
		if err == db.ErrNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			data, _ := json.Marshal(models.NewErrorMessage())
			ctx.Write(data)
			return nil
		}
	}
	data, _ := json.Marshal(threadData)
	ctx.Write(data)
	return nil
}

func ThreadVote(ctx *routing.Context) error { return nil }
