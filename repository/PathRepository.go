package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func CreatePath(newPath models.Path) (interface{}, error) {
	// Insert the new user into the collection
	res, err := DB.Collection("paths").InsertOne(context.TODO(), newPath)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func FindPath(filter bson.M) (path models.Path, errorType int, err error) {
	err = DB.Collection("paths").FindOne(context.TODO(), filter).Decode(&path)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return path, constants.NoDocumentFound, nil
		} else {
			return path, constants.Other, err
		}
	}
	return path, constants.NoError, nil
}

func FindPaths(filter bson.M) (paths []models.Path, err error) {
	collection := DB.Collection("paths")
	// Find paths based on the filter
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err // Return an error if the find operation fails
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var path models.Path
		if err := cursor.Decode(&path); err != nil {
			return nil, err // Return an error if decoding fails
		}
		paths = append(paths, path)
	}
	if err := cursor.Err(); err != nil {
		return nil, err // Return an error if there was an issue during the cursor iteration
	}

	// Return the paths and the total count (you can adjust this based on your needs)
	return paths, nil
}
