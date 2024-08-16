package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func FindBookmark(conditions map[string]interface{}) (bookmark models.Bookmark, errorType int, err error) {
	filter := bson.M(conditions)

	err = DB.Collection("bookmarks").FindOne(context.TODO(), filter).Decode(&bookmark)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bookmark, constants.NoDocumentFound, nil
		} else {
			return bookmark, constants.Other, err
		}
	}
	return bookmark, constants.NoError, nil
}

func CreateBookmark(bookmark models.Bookmark) (interface{}, error) {
	// Insert the new user into the collection
	res, err := DB.Collection("bookmarks").InsertOne(context.TODO(), bookmark)
	if err != nil {

		return nil, err
	}
	return res.InsertedID, nil
}
