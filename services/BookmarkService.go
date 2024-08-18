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
	"strconv"
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

	bookmark.KeepReading = false

	filter = bson.M{"_id": bookmarkModel.MangaId}

	//Retreive the manga to check if there are new chapters to read
	mangaModel, code, err := repository.FindManga(filter)
	if err != nil {
		fmt.Println("Error obtaining manga:", err)
		return dtos.Bookmark{}, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	// Handle case where the bookmark was not found
	if code == constants.NoDocumentFound {
		return dtos.Bookmark{}, errors.New("El manga no existe")
	}

	if bookmarkModel.Status == constants.Reading {
		keepReading, err := compareNumbersInStrings(bookmarkModel.Chapter, mangaModel.TotalChapters)
		if err != nil {
			fmt.Println("Error converting and comparing manga and bookmark chapters:", err)
			fmt.Println("bookmark chapter:", bookmarkModel.Chapter)
			fmt.Println("manga chapter:", mangaModel.TotalChapters)
			return dtos.Bookmark{}, errors.New("Error interno")
		}

		bookmark.KeepReading = keepReading
	}

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
	filter := bson.M{"_id": objectID}

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

	//TODO obtain manga to put correct value of keepReading flag

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
func createNewBookmark(data dtos.CreateBookmark, mangaID, userID primitive.ObjectID) (string, error) {

	bookmark := models.Bookmark{
		UserId:   userID,
		MangaId:  mangaID,
		Chapter:  data.Chapter,
		Status:   data.Status,
		LastRead: primitive.NewDateTimeFromTime(time.Now()),
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

// Helper function to extractNumber extracts the first floating-point number from a string
func extractNumber(s string) (float64, error) {
	// Regular expression to match numbers, including decimals
	re := regexp.MustCompile(`-?\d+(\.\d+)?`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("no number found in string")
	}

	// Convert the matched string to a float
	number, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

// Helper function to compareNumbersInStrings compares numbers in two strings
func compareNumbersInStrings(str1, str2 string) (bool, error) {
	// Extract numbers from both strings
	num1, err := extractNumber(str1)
	if err != nil {
		return false, fmt.Errorf("error extracting number from first string: %v", err)
	}

	num2, err := extractNumber(str2)
	if err != nil {
		return false, fmt.Errorf("error extracting number from second string: %v", err)
	}

	// Compare the numbers
	return num1 < num2, nil
}
