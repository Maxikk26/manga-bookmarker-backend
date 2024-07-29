package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"manga-bookmarker-backend/models"
)

func GetUsers() (users []models.User, err error) {
	query, err := DB.Collection("user").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	// Unpacks the cursor into a slice
	if err = query.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}
