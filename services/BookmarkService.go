package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/repository"
	"strings"
	"time"
)

func CreateBookmark(data dtos.CreateBookmark) (bookmarkId string, err error) {

	// Find the index of "manga-"
	prefix := "manga-"
	idx := strings.Index(data.Url, prefix)
	if idx == -1 {
		fmt.Println("Prefix not found: ", data.Url)
		return bookmarkId, errors.New("No se encontro el prefijo de manganato en la url")
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
		return bookmarkId, errors.New("No se encontró el manga especificado")
	}

	if errorType == constants.NoDocumentFound {
		//Channel to push the result of the scrapper to a struct, so we can do next steps after we have the data
		ch := make(chan dtos.MangaScrapperData)

		//Call scraooer service as a goroutine
		go ScrapperService(data.Url, ch)

		mangaData := <-ch

		// Use dto-mapper to map the data to the struct
		err = mapper.Map(&manga, &mangaData)
		if err != nil {
			fmt.Println("Error mapping data:", err)
			return bookmarkId, err
		}

		//Set de manga identifier
		manga.Identifier = mangaIdentifier

		//Create manga
		id, err := repository.CreateManga(manga)
		if err != nil {
			fmt.Println("Error creating manga: ", err.Error())
			return bookmarkId, errors.New("Ocurrio un error creando el manga")
		}

		//parse id with type interface{} to primitive.ObjectID
		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			fmt.Println("id is not of type primitive.ObjectID")
			return bookmarkId, errors.New("Ocurrio un error creando el manga")
		}

		manga.Id = objectID
	}

	userID, err := primitive.ObjectIDFromHex(data.UserId)
	if err != nil {
		fmt.Println("Invalid ObjectID string:", err)
		return bookmarkId, errors.New("Error Interno")
	}

	conditions := map[string]interface{}{
		"manga_id": manga.Id,
		"user_id":  userID,
	}

	bookmark, code, err := repository.FindBookmark(conditions)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return bookmarkId, errors.New("Ocurrio un error obteniendo el marcador")
	}

	if code == constants.NoDocumentFound {
		fmt.Println(fmt.Sprintf("%+v", bookmark))
		bookmark.UserId = userID
		bookmark.MangaId = manga.Id
		bookmark.Chapter = data.Chapter
		bookmark.LastRead = primitive.NewDateTimeFromTime(time.Now())

		id, err := repository.CreateBookmark(bookmark)
		if err != nil {
			fmt.Println("Error creating bookmark:", err.Error())
			return bookmarkId, errors.New("Ocurrio un error creando el marcador")
		}

		fmt.Println("id", id)

		//TODO Fix conversion of the id that the create returns
		idString, ok := id.(string)
		if !ok {
			fmt.Println("The interface does not contain a string")
		}

		return idString, nil

	} else {
		return bookmarkId, errors.New("El marcador ya existe")
	}
}
