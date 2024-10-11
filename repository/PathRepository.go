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
