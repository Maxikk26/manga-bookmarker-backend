package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
	"time"
)

func CreateSiteConfig(siteConfig dtos.CreateSiteConfig) error {

	filter := bson.M{"name": siteConfig.Name}
	siteModel, _, err := repository.FindSiteConfig(filter)
	if err != nil {
		fmt.Println("Error obtaining site configuration: ", err.Error())
		return err
	}

	if siteModel.Name == siteConfig.Name {
		fmt.Println("Site configuration already exists")
		return errors.New("Site configuration already exists")
	}

	// Use dto-mapper to map the data to the struct
	err = utils.Mapper.Map(&siteModel, &siteConfig)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return err
	}

	siteModel.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	err = repository.CreateSiteConfig(siteModel)
	if err != nil {
		return err
	}

	return nil
}
