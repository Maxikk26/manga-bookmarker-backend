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
	err = DB.Collection("mangas").FindOne(context.TODO(), bson.M{key: value}).Decode(&manga)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return manga, constants.NoDocumentFound, nil
		} else {
			return manga, constants.Other, err
		}
	}
	return manga, constants.NoError, nil
}

func FindManga(filter bson.M) (manga models.Manga, errorType int, err error) {
	err = DB.Collection("mangas").FindOne(context.TODO(), filter).Decode(&manga)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return manga, constants.NoDocumentFound, nil
		} else {
			return manga, constants.Other, err
		}
	}
	return manga, constants.NoError, nil
}

func AllMangas() (mangas []models.Manga, code int, err error) {
	cursor, err := DB.Collection("mangas").Find(context.TODO(), bson.D{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, constants.NoDocumentFound, nil
		}
		return nil, constants.Other, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var manga models.Manga
		if err := cursor.Decode(&manga); err != nil {
			return nil, constants.Other, err
		}
		mangas = append(mangas, manga)
	}

	if err := cursor.Err(); err != nil {
		return nil, constants.Other, err
	}

	return mangas, constants.NoError, nil
}

func UpdateManga(filter bson.M, updates bson.D) (int, error) {
	result, err := DB.Collection("mangas").UpdateOne(context.TODO(), filter, updates)
	if err != nil {
		return constants.Other, err
	}

	if result.MatchedCount == 0 {
		return constants.NoDocumentFound, nil
	}

	return constants.NoError, nil

}
