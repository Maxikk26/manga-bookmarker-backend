package controllers

import (
	"github.com/kataras/iris/v12"
	"log"
	"manga-bookmarker-backend/services"
	"sync"
)

type CreateBookmarkRequest struct {
	Url             string `json:"url"`
	LastChapterRead string `json:"lastChapterRead,omitempty"`
}

func CreateBookmarkHandler(ctx iris.Context) {
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
	go services.CrawlerService(request.Url, &wg)
	wg.Wait()
	response.Ok = true
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
}
