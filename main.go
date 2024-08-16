package main

import (
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"log"
	"manga-bookmarker-backend/controllers"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
	"os"
	"time"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	//Load .env
	loadEnv()

	//Connect to DB.
	repository.Init()

	//Set timezone
	time.Local, _ = time.LoadLocation("America/Caracas")

	//Add convertion functions to mapper
	utils.AddConvertionFunctions()

	//Start iris server
	app := iris.New()
	api := app.Party("/api")
	{
		v1 := api.Party("/v1")
		{
			v1.Post("/login", controllers.LoginController)

			user := v1.Party("/users")
			{
				user.Get("", controllers.GetUsersHandler)
				user.Post("", controllers.CreateUserHandler)
			}
			bookmark := v1.Party("/bookmarks")
			{
				bookmark.Post("", controllers.CreateBookmarkHandler)
				bookmark.Get("/{id}", controllers.GetBookmarkHandler)
				bookmark.Get("", controllers.GetBookmarksHandler)
				bookmark.Patch("/{id}", controllers.UpdateBookmarkHandler)
			}
		}
	}
	host := ":" + os.Getenv("PORT")
	app.Listen(host)

}
