package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func CreateManga(data models.Manga) (id interface{}, err error) {
	// Insert the new user into the collection
	res, err := DB.Collection("mangas").InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func FindMangaByAny(key, value string) (manga models.Manga, errorType int, err error) {
	err = DB.Collection("users").FindOne(context.TODO(), bson.M{key: value}).Decode(&manga)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return manga, constants.NoDocumentFound, nil
		} else {
			return manga, constants.Other, err
		}
	}
	return manga, constants.NoError, nil
}
