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
	"os"
	"regexp"
	"time"
)

//Core services

func CreateBookmark(data dtos.CreateBookmark) (string, error) {
	const prefix = "manga-"

	// Extract manga identifier from URL
	mangaIdentifier, err := ExtractMangaIdentifier(data.Url, prefix)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("No se encontro el prefijo de manganato en la url")
	}

	// Find or scrape manga
	manga, err := FindOrScrapeManga(mangaIdentifier, data.Url)
	if err != nil {
		fmt.Println("Error obtaining manga:", err.Error())
		return "", errors.New("No se encontró el manga especificado")
	}

	// Convert UserId to ObjectID
	userID, err := primitive.ObjectIDFromHex(data.UserId)
	if err != nil {
		fmt.Println("Invalid ObjectID string:", err)
		return "", errors.New("Error Interno")
	}

	// Check if bookmark already exists
	existingBookmark, err := findExistingBookmark(manga.Id, userID)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return "", err
	}

	if existingBookmark != nil {
		return "", errors.New("El bookmark ya existe")
	}

	// Create and return new bookmark
	bookmarkId, err := createNewBookmark(data, manga.Id, userID)
	if err != nil {
		fmt.Println("Error creating bookmark:", err.Error())
		return "", errors.New("Ocurrio un error creando el bookmark")
	}

	return bookmarkId, nil
}

func CreateBookmarkV2(data dtos.CreateBookmark) (string, error) {
	siteId, _ := primitive.ObjectIDFromHex(data.SiteId)
	filter := bson.M{"_id": siteId}
	siteModel, code, err := repository.FindSiteConfig(filter)
	if err != nil || code == constants.NoDocumentFound {
		return "", errors.New("No se encontro la configuracion del sitio")
	}

	filter = bson.M{"path": data.Path}
	pathModel, code, err := repository.FindPath(filter)
	if err != nil {
		return "", errors.New("No se encontro la path")
	}

	if code == constants.NoDocumentFound {
		//Scrap manga from the site with its configurations
		ch := make(chan dtos.MangaScrapperData)
		go MangaScrappingV2(data.Path, siteModel, ch)

		mangaData := <-ch

		//Create Manga
		var manga models.Manga
		err = utils.Mapper.Map(&manga, &mangaData)
		if err != nil {
			fmt.Println(err)
			return "", errors.New("Error interno")
		}

		manga.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
		id, err := repository.CreateManga(manga)
		if err != nil {
			return "", errors.New("Error al crear el manga")
		}

		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			return "", errors.New("Ocurrio un error creando el manga")
		}

		//Create Path for the manga
		newPath := models.Path{
			MangaId:       objectID,
			SiteId:        siteId,
			Path:          data.Path,
			TotalChapters: mangaData.TotalChapters,
			LastUpdate:    primitive.NewDateTimeFromTime(mangaData.LastUpdate),
		}

		id, err = repository.CreatePath(newPath)
		if err != nil {
			fmt.Println("Error creating path: ", err.Error())
			return "", errors.New("Error al crear el path")
		}
		pathId, _ := id.(primitive.ObjectID)

		pathModel = newPath
		pathModel.Id = pathId
	}

	// Convert UserId to ObjectID
	userID, _ := primitive.ObjectIDFromHex(data.UserId)

	// Check if bookmark already exists
	existingBookmark, err := findExistingBookmark(pathModel.Id, userID)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return "", err
	}

	if existingBookmark != nil {
		return "", errors.New("El bookmark ya existe")
	}

	// Create and return new bookmark
	bookmarkId, err := createNewBookmark(data, pathModel.Id, userID)
	if err != nil {
		fmt.Println("Error creating bookmark:", err.Error())
		return "", errors.New("Ocurrio un error creando el bookmark")
	}

	return bookmarkId, nil
}

func BookmarkDetails(bookmarkId string) (dtos.Bookmark, error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bookmarkId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return dtos.Bookmark{}, errors.New("Id del bookmark inválido")
	}

	// Define conditions for finding the bookmark
	filter := bson.M{"_id": objectID}

	// Retrieve the bookmark from the repository
	bookmarkModel, code, err := repository.FindBookmark(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err)
		return dtos.Bookmark{}, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	// Handle case where the bookmark was not found
	if code == constants.NoDocumentFound {
		return dtos.Bookmark{}, errors.New("El bookmark no existe")
	}

	// Map the data from the model to the DTO
	var bookmark dtos.Bookmark
	if err := utils.Mapper.Map(&bookmark, &bookmarkModel); err != nil {
		fmt.Println("Error mapping data:", err)
		return dtos.Bookmark{}, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	bookmark.KeepReading = validateKeepReading(&bookmark)

	return bookmark, nil
}

func UserBookmarks(userId string) ([]dtos.Bookmark, error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return nil, errors.New("Id del usuario inválido")
	}

	// Define conditions for finding the user's bookmarks
	filter := bson.M{"userId": objectID}

	//TODO paginated search

	// Retrieve the bookmarks from the repository
	bookmarkModels, code, err := repository.FindBookmarks(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmarks:", err)
		return nil, errors.New("Ocurrió un error obteniendo los bookmarks")
	}

	// Handle case where no bookmarks were found
	if code == constants.NoDocumentFound {
		return nil, errors.New("El usuario no posee ningún bookmark")
	}

	// Map the data from the model to the DTOs
	var bookmarks []dtos.Bookmark
	if err := utils.Mapper.Map(&bookmarks, &bookmarkModels); err != nil {
		fmt.Println("Error mapping data:", err)
		return nil, errors.New("Error interno")
	}

	for i := range bookmarks {
		bookmarks[i].KeepReading = validateKeepReading(&bookmarks[i])
	}

	return bookmarks, nil
}

func UpdateBookmark(bookmarkId string, bookmark dtos.Bookmark) (dtos.Bookmark, error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bookmarkId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return dtos.Bookmark{}, errors.New("Id del bookmark inválido")
	}

	// Define conditions for finding the bookmark
	filter := bson.M{"_id": objectID}

	// Retrieve the existing bookmark from the repository
	existingBookmark, code, err := repository.FindBookmark(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return dtos.Bookmark{}, errors.New("Ocurrió un error obteniendo el bookmark")
	}

	// Handle case where the bookmark was not found
	if code == constants.NoDocumentFound {
		return dtos.Bookmark{}, errors.New("El bookmark especificado no existe")
	}

	// Update LastRead if the chapter has changed
	if existingBookmark.Chapter != bookmark.Chapter {
		bookmark.LastRead = time.Now().UTC()
	}

	//Map to another DTO to separate which values you can pass to the DB
	var bookmarkUpdate dtos.BookmarkUpdate
	err = utils.Mapper.Map(&bookmarkUpdate, &bookmark)
	if err != nil {
		fmt.Println("Error mapping data:", err)
		return dtos.Bookmark{}, errors.New("Error interno")
	}

	updateDoc, err := bson.Marshal(&bookmarkUpdate)
	if err != nil {
		fmt.Println("Error marshalling manga data:", err)
		return dtos.Bookmark{}, err
	}

	// Unmarshal the BSON into a bson.M map
	var updateBson bson.M
	err = bson.Unmarshal(updateDoc, &updateBson)
	if err != nil {
		fmt.Println("Error unmarshalling BSON:", err)
		return dtos.Bookmark{}, err
	}

	// Create the update document with $set
	update := bson.D{{"$set", updateBson}}
	update = append(update, bson.E{Key: "updatedAt", Value: primitive.NewDateTimeFromTime(time.Now())})

	// Update the bookmark in the repository
	code, err = repository.UpdateBookmark(filter, update)
	if err != nil {
		fmt.Println("Error updating bookmark:", err.Error())
		return dtos.Bookmark{}, errors.New("Ocurrió un error actualizando el bookmark")
	}

	// Retrieve the updated bookmark from the repository
	bookmarkModel, code, err := repository.FindBookmark(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return dtos.Bookmark{}, errors.New("Ocurrió un error obteniendo el bookmark")
	}

	// Map the updated data to the DTO
	var updatedBookmark dtos.Bookmark
	if err := utils.Mapper.Map(&updatedBookmark, &bookmarkModel); err != nil {
		fmt.Println("Error mapping data:", err)
		return dtos.Bookmark{}, errors.New("Error interno")
	}

	//Set the flag for keepReading in case the user is behind in chapters
	updatedBookmark.KeepReading = false
	/*mangaId, _ := primitive.ObjectIDFromHex(updatedBookmark.MangaId)

	filter = bson.M{"_id": mangaId}
	mangaModel, code, err := repository.FindManga(filter)
	if err == nil && code == constants.NoError {
		updatedBookmark.KeepReading = updatedBookmark.Chapter < mangaModel.TotalChapters
	}*/

	return updatedBookmark, nil
}

func CheckForMangaUpdates(bookmarkId string) (dtos.Bookmark, error) {
	//obtain bookmark

	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bookmarkId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return dtos.Bookmark{}, errors.New("Id del bookmark inválido")
	}

	// Define conditions for finding the bookmark
	filter := bson.M{"_id": objectID}

	existingBookmark, code, err := repository.FindBookmark(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return dtos.Bookmark{}, err
	}
	if code == constants.NoDocumentFound {
		return dtos.Bookmark{}, nil
	}

	//obtain manga

	filter = bson.M{"_id": existingBookmark.MangaId}

	existingManga, code, err := repository.FindManga(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return dtos.Bookmark{}, err
	}

	if code == constants.NoDocumentFound {
		return dtos.Bookmark{}, nil
	}

	//web scrapping only looking for last chapter and date of udpdate

	url := os.Getenv("MANGANATO_URL") + existingManga.Identifier

	ch := make(chan dtos.MangaScrapperData)
	go SyncUpdatesScrapping(url, ch)
	mangaData := <-ch

	//update manga if there are any changes
	var bookmark dtos.Bookmark
	_ = utils.Mapper.Map(&bookmark, &existingBookmark)

	if existingManga.TotalChapters == mangaData.TotalChapters {
		if bookmark.Chapter != existingManga.TotalChapters {
			bookmark.KeepReading = true
		}
		return bookmark, nil
	}

	var updateMangaValues models.Manga
	err = utils.Mapper.Map(&updateMangaValues, &mangaData)
	if err != nil {
		return dtos.Bookmark{}, fmt.Errorf("Error mapping data: %v", err)
	}

	updateDoc, err := bson.Marshal(&mangaData)
	if err != nil {
		fmt.Println("Error marshalling manga data:", err)
		return dtos.Bookmark{}, err
	}

	// Unmarshal the BSON into a bson.M map
	var updateBson bson.M
	err = bson.Unmarshal(updateDoc, &updateBson)
	if err != nil {
		fmt.Println("Error unmarshalling BSON:", err)
		return dtos.Bookmark{}, err
	}

	// Create the update document with $set
	update := bson.D{{"$set", updateBson}}
	update = append(update, bson.E{Key: "updatedAt", Value: primitive.NewDateTimeFromTime(time.Now())})

	code, err = repository.UpdateManga(filter, update)
	if err != nil {
		return dtos.Bookmark{}, err
	}

	//return bookmark obj with keepReading and manga last chapter and date of update

	if mangaData.TotalChapters != bookmark.Chapter {
		bookmark.MangaUpdate = &dtos.BookmarkMangaUpdate{
			Update:      true,
			LastChapter: mangaData.TotalChapters,
			LastUpdated: &mangaData.LastUpdate,
		}
		bookmark.KeepReading = true
	}

	return bookmark, nil
}

//Helpers

// Helper function to find existing bookmark
func findExistingBookmark(mangaID, userID primitive.ObjectID) (*models.Bookmark, error) {
	filter := bson.M{
		"mangaId": mangaID,
		"userId":  userID,
	}

	bookmark, code, err := repository.FindBookmark(filter)
	if err != nil {
		return nil, err
	}

	if code == constants.NoDocumentFound {
		return nil, nil
	}

	return &bookmark, nil
}

// Helper function to create a new bookmark
func createNewBookmark(data dtos.CreateBookmark, pathId, userId primitive.ObjectID) (string, error) {

	bookmark := models.Bookmark{
		UserId:    userId,
		PathId:    pathId,
		Chapter:   data.Chapter,
		Status:    data.Status,
		LastRead:  primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	id, err := repository.CreateBookmark(bookmark)
	if err != nil {
		return "", err
	}

	objectID, ok := id.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Error parsing to the id of bookmark to ObjectID")
	}

	return objectID.Hex(), nil
}

// Helper function to validate the keepReading flag of a bookmark
func validateKeepReading(bookmark *dtos.Bookmark) bool {
	keepReading := false

	mangaId, _ := primitive.ObjectIDFromHex(bookmark.MangaId)
	filter := bson.M{"_id": mangaId}
	//Retreive the manga to check if there are new chapters to read
	mangaModel, code, err := repository.FindManga(filter)
	if err != nil {
		fmt.Println("Error obtaining manga:", err)
		return false
	}

	if code == constants.NoDocumentFound {
		fmt.Println("Manga does not exists", err)
		return false
	}

	keepReading = bookmark.Chapter < mangaModel.TotalChapters

	return keepReading
}

func extractPath(url string) (string, error) {
	// Regular expression to match and capture everything after the top-level domain
	re := regexp.MustCompile(`\.[a-z]{2,3}\/(.+)$`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		return "/" + match[1], nil
	}
	return "", fmt.Errorf("no match found")
}
