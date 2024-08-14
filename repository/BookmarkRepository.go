package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func FindBookmarkByAny(key, value string) (bookmark models.Bookmark, errorType int, err error) {
	err = DB.Collection("bookmarks").FindOne(context.TODO(), bson.M{key: value}).Decode(&bookmark)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return bookmark, constants.NoDocumentFound, nil
		} else {
			return bookmark, constants.Other, err
		}
	}
	return bookmark, constants.NoError, nil
}
