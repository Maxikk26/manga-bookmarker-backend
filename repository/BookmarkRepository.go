package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func FindBookmarksV2(filter bson.M, pageSize int, firstId primitive.ObjectID, lastId primitive.ObjectID) ([]models.Bookmark, int, error) {
	// Ensure that the pageSize is valid (greater than 0)
	if pageSize <= 0 {
		return nil, 0, errors.New("invalid page size")
	}

	collection := DB.Collection("bookmarks")

	// Specify the options for the query
	findOptions := options.Find()
	findOptions.SetLimit(int64(pageSize)) // Limit the number of bookmarks fetched

	// Adjust the filter and sort options based on firstId and lastId
	if !lastId.IsZero() {
		// Forward Pagination: Get the next set of bookmarks (lastId -> ...)
		filter["_id"] = bson.M{"$gt": lastId}               // Greater than lastId
		findOptions.SetSort(bson.D{{Key: "_id", Value: 1}}) // Ascending order by _id
	} else if !firstId.IsZero() {
		// Backward Pagination: Get the previous set of bookmarks (< firstId)
		filter["_id"] = bson.M{"$lt": firstId}               // Less than firstId
		findOptions.SetSort(bson.D{{Key: "_id", Value: -1}}) // Descending order by _id
	} else {
		// Default behavior (forward pagination from the beginning)
		findOptions.SetSort(bson.D{{Key: "_id", Value: 1}}) // Ascending order by _id
	}

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

	// If using backward pagination (firstId), reverse the result set
	if !firstId.IsZero() {
		for i, j := 0, len(bookmarks)-1; i < j; i, j = i+1, j-1 {
			bookmarks[i], bookmarks[j] = bookmarks[j], bookmarks[i]
		}
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
