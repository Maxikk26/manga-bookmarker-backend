package repository

import (
	"context"
	"errors"
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

func ListAllSiteConfigs() (siteConfigs []models.SiteConfig, errorType int, err error) {
	cursor, err := DB.Collection("siteConfigs").Find(context.TODO(), bson.D{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, constants.NoDocumentFound, nil
		}
		return nil, constants.Other, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var siteConfig models.SiteConfig
		if err := cursor.Decode(&siteConfig); err != nil {
			return nil, constants.Other, err
		}
		siteConfigs = append(siteConfigs, siteConfig)
	}

	if err := cursor.Err(); err != nil {
		return nil, constants.Other, err
	}

	return siteConfigs, constants.NoError, nil
}
