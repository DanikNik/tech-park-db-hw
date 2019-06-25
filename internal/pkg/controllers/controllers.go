package controllers

import (
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/models"
)

func Index(ctx *routing.Context) error { return nil }

func Clear(ctx *routing.Context) error {
	db.Truncate()
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func Status(ctx *routing.Context) error {
	status := &models.Status{}
	db.Status(status)
	ctx.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(status)
	ctx.Write(data)
	return nil
}
