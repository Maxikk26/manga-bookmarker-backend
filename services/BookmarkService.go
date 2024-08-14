package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/repository"
	"strings"
)

func CreateBookmark(data dtos.CreateBookmark) error {

	// Find the index of "manga-"
	prefix := "manga-"
	idx := strings.Index(data.Url, prefix)
	if idx == -1 {
		fmt.Println("Prefix not found: ", data.Url)
		return errors.New("No se encontro el prefijo de manganato en la url")
	}

	// Extract the substring after "manga-"
	mangaIdentifier := data.Url[idx+len(prefix):]

	// Check if there's a "/" in the extracted ID and handle it accordingly
	if slashIdx := strings.Index(mangaIdentifier, "/"); slashIdx != -1 {
		mangaIdentifier = mangaIdentifier[:slashIdx] // Exclude the part after "/"
	}

	fmt.Println("mangaIdentifier", mangaIdentifier)

	manga, errorType, err := repository.FindMangaByAny("identifier", mangaIdentifier)
	if err != nil {
		fmt.Println("Error obtaining manga:", err.Error())
		return errors.New("No se encontró el manga especificado")
	}

	if errorType == constants.NoDocumentFound {
		//TODO scrap manga info to poblate manga struct
		//Channel to push the result of the scrapper to a struct so we can do nexts steps after we have the data
		ch := make(chan dtos.MangaScrapperData)

		//Call scraooer service as a goroutine
		go ScrapperService(data.Url, ch)

		mangaData := <-ch

		// Use dto-mapper to map the data to the struct
		err = mapper.Map(&manga, &mangaData)
		if err != nil {
			fmt.Println("Error mapping data:", err)
			return err
		}

		//Set de manga identifier
		manga.Identifier = mangaIdentifier

		//Create manga
		id, err := repository.CreateManga(manga)
		if err != nil {
			fmt.Println("Error creating manga: ", err.Error())
			return errors.New("Ocurrio un error creando el manga")
		}

		//parse id with type interface{} to primitive.ObjectID
		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			fmt.Println("id is not of type primitive.ObjectID")
			return errors.New("Ocurrio un error creando el manga")
		}

		manga.Id = objectID
	}

	//TODO check if the bookmark for the manga on current user does not exists

	//TODO create the bookmark

	return nil
}
