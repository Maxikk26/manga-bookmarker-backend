package controllers

import (
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/services"
)

func LoginController(ctx iris.Context) {
	var response Response

	var request dtos.Login
	if err := ctx.ReadJSON(&request); err != nil {
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	jwt, err := services.Login(request)
	if err != nil {
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(response)
		return
	}

	response.Ok = true
	response.Result = jwt
	ctx.JSON(response)
	return
}
