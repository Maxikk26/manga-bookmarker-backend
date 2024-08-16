package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func GetUsers() (users []models.User, err error) {
	query, err := DB.Collection("users").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	// Unpacks the cursor into a slice
	if err = query.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(newUser models.User) error {
	// Insert the new user into the collection
	_, err = DB.Collection("users").InsertOne(context.TODO(), newUser)
	if err != nil {
		return err
	}
	return nil
}

func FindUserByAny(key, value string) (user models.User, err error) {
	err = DB.Collection("users").FindOne(context.TODO(), bson.M{key: value}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("User not found")
		} else {
			return user, err
		}
	}
	return user, nil
}

func FindUser(conditions map[string]interface{}) (user models.User, errorType int, err error) {
	filter := bson.M(conditions)

	err = DB.Collection("bookmarks").FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user, constants.NoDocumentFound, nil
		} else {
			return user, constants.Other, err
		}
	}
	return user, constants.NoError, nil
}
