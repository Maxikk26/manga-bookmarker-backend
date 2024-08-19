package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/models"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
	"strings"
)

//Core functions

func AllMangas() (mangas []dtos.MangaInfo, err error) {

	mangaModel, code, err := repository.AllMangas()
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return mangas, errors.New("Ocurrio un error obteniendo los mangas")
	}

	if code == constants.NoDocumentFound {
		return mangas, errors.New("No existe ningún manga")
	}

	err = utils.Mapper.Map(&mangas, &mangaModel)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return mangas, errors.New("Error Interno")
	}

	return mangas, nil
}

//Helpers

// UpdateManga Helper to update a manga using the dto MangaScrapperData
func UpdateManga(data dtos.MangaScrapperData, filter bson.M) (err error) {
	var updateMangaValues models.Manga
	err = utils.Mapper.Map(&updateMangaValues, &data)
	if err != nil {
		fmt.Println(err)
	}

	updateDoc, err := bson.Marshal(&data)
	if err != nil {
		fmt.Println("Error marshalling manga data:", err)
	}

	// Unmarshal the BSON into a bson.M map
	var updateBson bson.M
	err = bson.Unmarshal(updateDoc, &updateBson)
	if err != nil {
		fmt.Println("Error unmarshalling BSON:", err)
	}

	// Create the update document with $set
	update := bson.D{{"$set", updateBson}}

	_, err = repository.UpdateManga(filter, update)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// FindOrScrapeManga Helper function to find or scrape manga
func FindOrScrapeManga(mangaIdentifier, url string) (models.Manga, error) {
	filter := bson.M{"identifier": mangaIdentifier}
	manga, errorType, err := repository.FindManga(filter)
	if err != nil {
		return models.Manga{}, err
	}

	if errorType == constants.NoDocumentFound {
		ch := make(chan dtos.MangaScrapperData)
		go MangaScrapping(url, ch)

		mangaData := <-ch
		err = utils.Mapper.Map(&manga, &mangaData)
		if err != nil {
			return models.Manga{}, fmt.Errorf("Error mapping data: %v", err)
		}

		manga.Identifier = mangaIdentifier
		id, err := repository.CreateManga(manga)
		if err != nil {
			return models.Manga{}, fmt.Errorf("Error creating manga: %v", err)
		}

		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			return models.Manga{}, errors.New("Ocurrio un error creando el manga")
		}

		manga.Id = objectID
	}

	return manga, nil
}

// ExtractMangaIdentifier Helper function to extract manga identifier from URL
func ExtractMangaIdentifier(url, prefix string) (string, error) {
	idx := strings.Index(url, prefix)
	if idx == -1 {
		return "", fmt.Errorf("Prefix not found: %s", url)
	}

	mangaIdentifier := url[idx+len(prefix):]
	if slashIdx := strings.Index(mangaIdentifier, "/"); slashIdx != -1 {
		mangaIdentifier = mangaIdentifier[:slashIdx]
	}

	return mangaIdentifier, nil
}
