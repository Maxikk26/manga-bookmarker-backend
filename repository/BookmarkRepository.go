package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func FindBookmark(filter bson.M) (bookmark models.Bookmark, errorType int, err error) {
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

func FindBookmarks(filter bson.M) (bookmarks []models.Bookmark, code int, err error) {
	cursor, err := DB.Collection("bookmarks").Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.NoDocumentFound, nil
		}
		return nil, constants.Other, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var bookmark models.Bookmark
		if err := cursor.Decode(&bookmark); err != nil {
			return nil, constants.Other, err
		}
		bookmarks = append(bookmarks, bookmark)
	}

	if err := cursor.Err(); err != nil {
		return nil, constants.Other, err
	}

	return bookmarks, constants.NoError, nil
}

func CreateBookmark(bookmark models.Bookmark) (interface{}, error) {
	// Insert the new user into the collection
	res, err := DB.Collection("bookmarks").InsertOne(context.TODO(), bookmark)
	if err != nil {

		return nil, err
	}
	return res.InsertedID, nil
}

func UpdateBookmark(filter bson.M, updates bson.D) (int, error) {
	result, err := DB.Collection("bookmarks").UpdateOne(context.TODO(), filter, updates)
	if err != nil {
		return constants.Other, err
	}

	if result.MatchedCount == 0 {
		return constants.NoDocumentFound, nil
	}

	return constants.NoError, nil

}
