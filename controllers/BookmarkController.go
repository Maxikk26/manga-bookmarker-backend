package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/services"
	"strings"
)

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

	id, err := services.CreateBookmark(request)
	if err != nil {
		fmt.Println("Error while creating bookmark: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		return
	}

	result := map[string]interface{}{
		"bookmarkId": id,
	}

	response.Ok = true
	response.Result = result
	ctx.JSON(response)
	return
}

func GetBookmarkHandler(ctx iris.Context) {
	var response Response

	bookmarkId := ctx.Params().Get("id")

	if bookmarkId == "" {
		response.Ok = false
		response.Msg = "El parametro id está vacío"
		ctx.JSON(response)
		return
	}

	result, err := services.BookmarkDetails(bookmarkId)
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

func GetBookmarksHandler(ctx iris.Context) {
	var response Response

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

	result, err := services.UserBookmarks(userId)
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

func UpdateBookmarkHandler(ctx iris.Context) {
	var response Response

	bookmarkId := ctx.Params().Get("id")

	if bookmarkId == "" {
		response.Ok = false
		response.Msg = "El parametro id está vacío"
		ctx.JSON(response)
		return
	}

	var request dtos.Bookmark
	if err := ctx.ReadJSON(&request); err != nil {
		fmt.Println("Error while parsing request body: ", err.Error())
		response.Ok = false
		response.Msg = err.Error()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(response)
		return
	}

	result, err := services.UpdateBookmark(bookmarkId, request)
	if err != nil {
		fmt.Println("Error while updating bookmark: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		return
	}

	response.Ok = true
	response.Result = result
	ctx.JSON(response)
}
