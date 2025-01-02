package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func FindBookmarksV2(filter bson.M, pageSize int) ([]models.Bookmark, int, error) {

	// Ensure that the pageSize is valid (greater than 0)
	if pageSize <= 0 {
		return nil, 0, errors.New("invalid page size")
	}

	collection := DB.Collection("bookmarks")

	// Specify the options for the query
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize)) // Limit the number of bookmarks fetched

	// Optionally, sort the results by _id to maintain order
	findOptions.SetSort(bson.D{{Key: "_id", Value: 1}}) // Ascending order by _id

	// Execute the query to find bookmarks
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("error finding bookmarks: %w", err)
	}
	defer cursor.Close(context.Background())

	// Slice to hold the resulting bookmarks
	var bookmarks []models.Bookmark
	for cursor.Next(context.Background()) {
		var bookmark models.Bookmark
		if err := cursor.Decode(&bookmark); err != nil {
			return nil, 0, fmt.Errorf("error decoding bookmark: %w", err)
		}
		bookmarks = append(bookmarks, bookmark)
	}

	// Check if there was an error during the iteration
	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	// Return the result with the number of bookmarks found
	return bookmarks, len(bookmarks), nil
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
