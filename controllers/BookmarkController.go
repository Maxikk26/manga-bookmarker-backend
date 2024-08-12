package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/services"
	"strings"
)

/*func CreateBookmarkHandler(ctx iris.Context) {
	var response Response

	var request CreateBookmarkRequest
	err := ctx.ReadJSON(&request)
	if err != nil {
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(response)
		log.Fatal(err)
		return
	}
	//Wait group to listen when the goroutine finished and get back at the user
	var wg sync.WaitGroup
	wg.Add(1)
	//Call Crawler service as a goroutine
	go services.ScrapperService(request.Url, &wg)
	wg.Wait()
	response.Ok = true
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
}*/

func CreateBookmarkHandler(ctx iris.Context) {
	var response Response

	var request dtos.CreateBookmark
	if err := ctx.ReadJSON(&request); err != nil {
		fmt.Println("Error while parsing request body: ", err.Error())
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	// Extract the Authorization header
	authHeader := ctx.GetHeader("Authorization")
	// Extract the token by trimming the "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")

	userId, err := services.GetUserIdFromClaims(token)
	if err != nil {
		fmt.Println("Error obtaining claims: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	request.UserId = userId

	if err = services.CreateBookmark(request); err != nil {
		fmt.Println("Error while creating bookmark: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		return
	}

	response.Ok = true
	ctx.JSON(response)
	return
}
