package controllers

import (
	"encoding/json"
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"tech-park-db-hw/internal/pkg/db"
	"tech-park-db-hw/internal/pkg/models"
)

func UserCreate(ctx *routing.Context) error {
	nick := ctx.Param("nickname")

	var userData models.User
	err := json.Unmarshal(ctx.PostBody(), &userData)
	if err != nil {
		_, _ = fmt.Fprintln(ctx, err)
		return err
	}
	userData.Nickname = nick

	err = db.CreateUser(&userData)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusConflict)
		userList, err := db.SelectUsersOnConflict(userData.Nickname, userData.Email)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return err
		}
		fmt.Println(userList)
		data, err := json.Marshal(&userList)
		ctx.Write(data)
		return err
	}
	data, err := userData.MarshalJSON()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write(data)
	return nil
}

func UserGetOne(ctx *routing.Context) error {
	//nick := ctx.Param("nickname")
	return nil
}

func UserUpdate(ctx *routing.Context) error {
	return nil
}
