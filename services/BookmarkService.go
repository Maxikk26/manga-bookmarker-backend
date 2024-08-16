package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
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
		err = utils.Mapper.Map(&manga, &mangaData)
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
		"mangaId": manga.Id,
		"userId":  userID,
	}

	bookmark, code, err := repository.FindBookmark(conditions)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return bookmarkId, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	if code == constants.NoDocumentFound {
		bookmark.UserId = userID
		bookmark.MangaId = manga.Id
		bookmark.Chapter = data.Chapter
		bookmark.Status = data.Status
		bookmark.LastRead = primitive.NewDateTimeFromTime(time.Now())

		id, err := repository.CreateBookmark(bookmark)
		if err != nil {
			fmt.Println("Error creating bookmark:", err.Error())
			return bookmarkId, errors.New("Ocurrio un error creando el bookmark")
		}

		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			fmt.Println("Error parsing to the id of bookmark to ObjectID")
			return bookmarkId, errors.New("Ocurrio un error enviando la respuesta")
		}

		return objectID.Hex(), nil

	} else {
		return bookmarkId, errors.New("El bookmark ya existe")
	}
}

func BookmarkDetails(bookmarkId string) (bookmark dtos.Bookmark, err error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bookmarkId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return bookmark, errors.New("Id del bookmark inválido")
	}

	conditions := map[string]interface{}{
		"_id": objectID,
	}

	bookmarkModel, code, err := repository.FindBookmark(conditions)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return bookmark, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	if code == constants.NoDocumentFound {
		return bookmark, errors.New("El bookmark no existe")
	}

	// Use dto-mapper to map the data to the struct
	err = utils.Mapper.Map(&bookmark, &bookmarkModel)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return bookmark, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	return bookmark, nil
}

func UserBookmarks(userId string) (bookmarks []dtos.Bookmark, err error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return bookmarks, errors.New("Id del bookmark inválido")
	}

	conditions := map[string]interface{}{
		"userId": objectID,
	}

	bookmarkModel, code, err := repository.FindBookmarks(conditions)
	if err != nil {
		fmt.Println("Error obtaining bookmarks:", err.Error())
		return bookmarks, errors.New("Ocurrio un error obteniendo los bookmarks")
	}

	if code == constants.NoDocumentFound {
		return bookmarks, errors.New("El usuario no posee ningún bookmark")
	}

	err = utils.Mapper.Map(&bookmarks, &bookmarkModel)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return bookmarks, errors.New("Error Interno")
	}

	return bookmarks, nil
}

func UpdateBookmark(bookmarkId string, bookmark dtos.Bookmark) (dtos.Bookmark, error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bookmarkId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return bookmark, errors.New("Id del bookmark inválido")
	}

	conditions := map[string]interface{}{
		"_id": objectID,
	}

	updates := utils.StructToMap(bookmark)

	code, err := repository.UpdateBookmark(conditions, updates)
	if err != nil {
		fmt.Println("Error obtaining bookmarks:", err.Error())
		return bookmark, errors.New("Ocurrio un error actualizando el bookmark")
	}

	if code == constants.NoDocumentFound {
		return bookmark, errors.New("El bookmark especificado no existe")
	}

	//TODO obtain updated bookmark

	bookmarkModel, code, err := repository.FindBookmark(conditions)
	if err != nil {
		fmt.Println("Error obtaining bookmarks:", err.Error())
		return bookmark, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	// Use dto-mapper to map the data to the struct
	err = utils.Mapper.Map(&bookmark, &bookmarkModel)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return bookmark, errors.New("Error interno")
	}

	return bookmark, nil

}
