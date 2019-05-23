package controllers

import (
	"encoding/json"
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"tech-park-db-hw/internal/pkg/models"
)

func UserCreate(ctx *routing.Context) error {
	var userData models.User
	err := json.Unmarshal(ctx.PostBody(), &userData)
	if err != nil {
		_, _ = fmt.Fprintln(ctx, err)
		return err
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	_, _ = fmt.Fprint(ctx, ctx.Param("nickname"))
	return nil
}

func UserGetOne(ctx *routing.Context) error {
	return nil
}

func UserUpdate(ctx *routing.Context) error {
	return nil
}
