package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
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

	variable := os.Getenv("MIGRATIONS")
	migrations, err := strconv.ParseBool(variable)
	if err != nil {
		log.Fatalf("Error converting environment variable to boolean: %v\n", err)
	}

	if migrations {
		makeMigrations()
	}

	fmt.Println("Connected to MongoDB!")

}

func makeMigrations() {
	// List of collection names to create
	collections := []string{"users", "mangas", "bookmarks", "siteConfigs", "paths"}
	// Create each collection
	for _, collectionName := range collections {
		err = DB.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			log.Fatalf("Error creating collection '%s': %v\n", collectionName, err)
		}
		fmt.Printf("Collection '%s' created successfully\n", collectionName)
	}

}
