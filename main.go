package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/robfig/cron"
	"log"
	"manga-bookmarker-backend/controllers"
	"manga-bookmarker-backend/middleware"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/services"
	"manga-bookmarker-backend/utils"
	"os"
	"strconv"
	"time"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func loadScrapperCron() {
	variable := os.Getenv("SCRAPPER_CRON")
	scrapperFlag, err := strconv.ParseBool(variable)
	if err != nil {
		log.Fatalf("Error converting environment variable to boolean: %v\n", err)
	}

	if scrapperFlag {
		c := cron.New()

		err = c.AddFunc(os.Getenv("SCRAPPER_CRON_SCHEDULE"), services.ScrappingJob)
		if err != nil {
			fmt.Println(err)
			return
		}

		c.Start()

	}
	return
}

func main() {
	//Set timezone
	time.Local, _ = time.LoadLocation("America/Caracas")

	//Load .env
	loadEnv()

	//Connect to DB.
	repository.Init()

	//Set up scrapper cron
	//loadScrapperCron()

	//Add convertion functions to mapper
	utils.AddConvertionFunctions()

	//Start iris server
	app := iris.New()
	api := app.Party("/api")
	{
		v1 := api.Party("/v1")
		{
			v1.Post("/login", controllers.LoginController)

			user := v1.Party("/users", middleware.Auth)
			{
				user.Get("", controllers.GetUsersHandler)
				user.Post("", controllers.CreateUserHandler)
			}
			bookmark := v1.Party("/bookmarks", middleware.Auth)
			{
				bookmark.Post("", controllers.CreateBookmarkHandler)
				bookmark.Get("/{id}", controllers.GetBookmarkHandler)
				bookmark.Get("", controllers.GetBookmarksHandler)
				bookmark.Get("/{id}/manga", controllers.CheckUpdatesHandler)
				bookmark.Patch("/{id}", controllers.UpdateBookmarkHandler)
			}
			manga := v1.Party("/mangas", middleware.Auth)
			{
				manga.Get("", controllers.GetMangasHandler)
			}
			site := v1.Party("/sites", middleware.Auth)
			{
				site.Post("", controllers.CreateSiteConfigHandler)
				site.Get("/selector", controllers.ListSiteConfigHandler)
			}
		}
	}
	host := ":" + os.Getenv("PORT")
	app.Listen(host)

}
