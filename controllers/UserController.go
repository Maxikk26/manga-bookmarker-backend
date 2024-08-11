package controllers

import (
	"github.com/dranikpg/dto-mapper"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/services"
)

type Response struct {
	Ok     bool        `json:"ok"`
	Msg    string      `json:"msg,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

var mapper dto.Mapper

func GetUsersHandler(ctx iris.Context) {
	var response Response

	users, err := services.GetAllUsers()
	if err != nil {
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	response.Ok = true
	response.Result = users
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
	return
}

func CreateUserHandler(ctx iris.Context) {
	var response Response

	var request dtos.UserCreate
	if err := ctx.ReadJSON(&request); err != nil {
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	err := services.CreateUser(request)
	if err != nil {
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(response)
		return
	}

	response.Ok = true
	ctx.JSON(response)
	return

}
