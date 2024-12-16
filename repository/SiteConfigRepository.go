package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/models"
)

func CreateSiteConfig(newSiteConfig models.SiteConfig) error {
	// Insert the new user into the collection
	_, err = DB.Collection("siteConfigs").InsertOne(context.TODO(), newSiteConfig)
	if err != nil {
		return err
	}
	return nil
}

func FindSiteConfig(filter bson.M) (siteConfig models.SiteConfig, errorType int, err error) {
	err = DB.Collection("siteConfigs").FindOne(context.TODO(), filter).Decode(&siteConfig)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return siteConfig, constants.NoDocumentFound, nil
		} else {
			return siteConfig, constants.Other, err
		}
	}
	return siteConfig, constants.NoError, nil
}
