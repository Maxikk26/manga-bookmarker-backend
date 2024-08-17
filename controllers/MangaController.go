package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/services"
)

func GetMangasHandler(ctx iris.Context) {
	var response Response

	result, err := services.AllMangas()
	if err != nil {
		fmt.Println("Error while getting bookmark detail: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	response.Ok = true
	response.Result = result
	ctx.JSON(response)
	return
}
