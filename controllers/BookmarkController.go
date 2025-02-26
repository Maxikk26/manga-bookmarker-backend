package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/services"
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

	request.UserId = ctx.Values().Get("userId").(string)

	//id, err := services.CreateBookmark(request)
	id, err := services.CreateBookmarkV2(request)
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

type BookmarkSearchParams struct {
	FirstId  string
	LastId   string
	PageSize int
	Status   int
}

func GetBookmarksHandler(ctx iris.Context) {
	var response Response

	userId := ctx.Values().Get("userId").(string)
	firstIdStr := ctx.URLParam("firstId")
	lastIdStr := ctx.URLParam("lastId")
	pageSize := ctx.URLParamIntDefault("pageSize", 5)

	result, err := services.UserBookmarks(userId, firstIdStr, lastIdStr, pageSize)
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

func CheckUpdatesHandler(ctx iris.Context) {
	var response Response

	bookmarkId := ctx.Params().Get("id")

	if bookmarkId == "" {
		response.Ok = false
		response.Msg = "El parametro id está vacío"
		ctx.JSON(response)
		return
	}

	//TODO refactor

	/*result, err := services.CheckForMangaUpdates(bookmarkId)
	if err != nil {
		fmt.Println("Error while checking for manga updates: ", err)
		response.Ok = false
		response.Msg = err.Error()
		ctx.JSON(response)
		return
	}*/

	response.Ok = true
	//response.Result = result
	ctx.JSON(response)

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
