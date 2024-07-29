package main

import (
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"log"
	"manga-bookmarker-backend/controllers"
	"manga-bookmarker-backend/repository"
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

	//Start iris server
	app := iris.New()
	api := app.Party("/api")
	{
		v1 := api.Party("/v1")
		{
			//API
			user := v1.Party("/users")
			{
				user.Get("", controllers.GetUsersHandler)
			}
			bookmark := v1.Party("/bookmarks")
			{
				bookmark.Post("", controllers.CreateBookmarkHandler)
			}
		}
	}
	app.Listen(":8080")

}
