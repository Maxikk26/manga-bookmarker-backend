package services

import (
	"log"
	"manga-bookmarker-backend/models"
	"manga-bookmarker-backend/repository"
)

func GetAllUsers() (users []models.User, err error) {
	users, err = repository.GetUsers()
	if err != nil {
		log.Fatal("GetAllUsers Error: ", err)
		return nil, err
	}
	return users, nil
}
