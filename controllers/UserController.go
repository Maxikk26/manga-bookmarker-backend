package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/services"
)

type Response struct {
	Ok     bool        `json:"ok"`
	Msg    string      `json:"msg,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

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
	fmt.Print("CreateUserHandler")
}
