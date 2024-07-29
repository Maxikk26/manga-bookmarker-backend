package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	DB  *mongo.Database
	err error
)

func Init() {
	//Connect DB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " +
			"www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	DB = client.Database("manga-bookmark")
	//defer func() {
	//	if err := client.Disconnect(context.TODO()); err != nil {
	//		panic(err)
	//	}
	//}()

	fmt.Println("Connected to MongoDB!")

}
