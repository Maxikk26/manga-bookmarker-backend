package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/services"
)

func CreateSiteConfigHandler(ctx iris.Context) {
	var response Response

	var request dtos.CreateSiteConfig
	if err := ctx.ReadJSON(&request); err != nil {
		fmt.Println("Error while parsing request body: ", err.Error())
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	err := services.CreateSiteConfig(request)
	if err != nil {
		fmt.Println("Error while creating site configuration: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		return
	}

	response.Ok = true
	ctx.JSON(response)
	return
}

func ListSiteConfigHandler(ctx iris.Context) {
	var response Response

	siteConfigs, err := services.ListSites()
	if err != nil {
		fmt.Println("Error while creating site configuration: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		return
	}

	response.Ok = true
	response.Result = siteConfigs
	ctx.JSON(response)
	return
}
