package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/models"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
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

	path := utils.PathFromURL(data.Url)

	filter = bson.M{"path": path}
	pathModel, code, err := repository.FindPath(filter)
	if err != nil {
		return "", errors.New("No se encontro la path")
	}

	if code == constants.NoDocumentFound {
		//Scrap manga from the site with its configurations
		ch := make(chan dtos.MangaScrapperData)
		go MangaScrappingV2(path, siteModel, ch)

		mangaData := <-ch
		log.Print("mangaData", mangaData)

		//Create Manga
		var manga models.Manga
		err = utils.Mapper.Map(&manga, &mangaData)
		if err != nil {
			fmt.Println(err)
			return "", errors.New("Error interno")
		}

		manga.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
		log.Println("manga model", manga)
		id, err := repository.CreateManga(manga)
		if err != nil {
			return "", errors.New("Error al crear el manga")
		}

		log.Println("manga id", id)

		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			return "", errors.New("Ocurrio un error creando el manga")
		}

		//Create Path for the manga
		newPath := models.Path{
			MangaId:       objectID,
			SiteId:        siteId,
			Path:          path,
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
	existingBookmark, err := findExistingBookmarkV2(pathModel.Id, userID)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err.Error())
		return "", err
	}
	fmt.Println("existingBookmark", existingBookmark)

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

func BookmarkDetails(bookmarkId string) (dtos.BookmarkDetail, error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(bookmarkId)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		return dtos.BookmarkDetail{}, errors.New("Id del bookmark inválido")
	}

	// Define conditions for finding the bookmark
	filter := bson.M{"_id": objectID}

	// Retrieve the bookmark from the repository
	bookmarkModel, code, err := repository.FindBookmark(filter)
	if err != nil {
		fmt.Println("Error obtaining bookmark:", err)
		return dtos.BookmarkDetail{}, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	// Handle case where the bookmark was not found
	if code == constants.NoDocumentFound {
		return dtos.BookmarkDetail{}, errors.New("El bookmark no existe")
	}

	// Map the data from the model to the DTO
	var bookmark dtos.BookmarkDetail
	if err := utils.Mapper.Map(&bookmark, &bookmarkModel); err != nil {
		fmt.Println("Error mapping data:", err)
		return dtos.BookmarkDetail{}, errors.New("Ocurrio un error obteniendo el bookmark")
	}

	filter = bson.M{"_id": bookmarkModel.PathId}
	// Retrieve the path from the DB
	pathModel, code, err := repository.FindPath(filter)
	if err != nil {
		fmt.Println("Error obtaining path:", err)
		return dtos.BookmarkDetail{}, errors.New("Ocurrio un error obteniendo el path del bookmark")
	}

	filter = bson.M{"_id": pathModel.MangaId}
	// Retrieve the path from the DB
	mangaModel, code, err := repository.FindManga(filter)
	if err != nil {
		fmt.Println("Error obtaining path:", err)
		return dtos.BookmarkDetail{}, errors.New("Ocurrio un error obteniendo el path del bookmark")
	}

	bookmark.MangaInfo.Name = mangaModel.Name
	bookmark.MangaInfo.TotalChapters = pathModel.TotalChapters
	bookmark.MangaInfo.LastUpdate = pathModel.LastUpdate.Time()

	return bookmark, nil
}

func UserBookmarks(userId, firstIdStr, lastIdStr string, pageSize int, count bool) (interface{}, error) {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(userId)

	// Define conditions for finding the user's bookmarks
	filter := bson.M{"userId": objectID}

	var firstId primitive.ObjectID
	if firstIdStr != "" {
		firstId, err = primitive.ObjectIDFromHex(firstIdStr)
		if err != nil {
			return nil, errors.New("El firstId enviado es invalido")
		}
	}

	var lastId primitive.ObjectID
	if lastIdStr != "" {
		lastId, err = primitive.ObjectIDFromHex(lastIdStr)
		if err != nil {
			return nil, errors.New("El lastId enviado es invalido")
		}
	}

	if !firstId.IsZero() && !lastId.IsZero() {
		return nil, errors.New("No se pueden usar firstId y lastId para paginacion, solo enviar 1.")
	}

	totalBookmarks := int64(0)

	//Retrieve total count of user bookmarks (only retrieving when using a flag)
	if count {
		totalBookmarks, _ = repository.CountUserBookmarks(filter)
	}

	// Retrieve the bookmarks from the repository
	bookmarkModels, code, err := repository.FindBookmarks(filter, pageSize, firstId, lastId)
	if err != nil {
		fmt.Println("Error obtaining bookmarks:", err)
		return nil, errors.New("Ocurrió un error obteniendo los bookmarks")
	}

	// Handle case where no bookmarks were found
	if code == constants.NoDocumentFound {
		return nil, errors.New("El usuario no posee ningún bookmark")
	}

	// Collect all PathIds from the bookmark models
	pathIds := make([]primitive.ObjectID, len(bookmarkModels))
	for i, bookmark := range bookmarkModels {
		pathIds[i] = bookmark.PathId
	}

	// Fetch all PathModels in a single query
	paths, err := repository.FindPaths(bson.M{"_id": bson.M{"$in": pathIds}})
	if err != nil {
		return []dtos.UserBookmars{}, errors.New("Error obteniendo paths")
	}

	// Map PathModels by PathId for easy lookup
	pathMap := make(map[primitive.ObjectID]models.Path)
	for _, path := range paths {
		pathMap[path.Id] = path
	}

	// Collect all MangaIds from the PathModels
	mangaIds := make([]primitive.ObjectID, 0, len(paths))
	for _, path := range paths {
		mangaIds = append(mangaIds, path.MangaId)
	}

	// Fetch all MangaModels in a single query
	mangas, err := repository.FindMangas(bson.M{"_id": bson.M{"$in": mangaIds}})
	if err != nil {
		return []dtos.UserBookmars{}, errors.New("Error fetching mangas")
	}

	// Map MangaModels by MangaId for easy lookup
	mangaMap := make(map[primitive.ObjectID]models.Manga)
	for _, manga := range mangas {
		mangaMap[manga.Id] = manga
	}

	var bookmarks []dtos.BookmarkDetail
	for _, bookmark := range bookmarkModels {
		pathModel, pathExists := pathMap[bookmark.PathId]
		if !pathExists {
			// Handle error if the path isn't found in the map
			continue
		}

		mangaModel, mangaExists := mangaMap[pathModel.MangaId]
		if !mangaExists {
			// Handle error if the manga isn't found in the map
			continue
		}

		// Determine if the user should keep reading
		keepReading := pathModel.TotalChapters > bookmark.Chapter

		// Build the detail of the bookmark
		bookmarkDetail := dtos.BookmarkDetail{
			Id:          bookmark.Id.Hex(),
			Chapter:     bookmark.Chapter,
			LastRead:    bookmark.LastRead.Time(),
			Status:      bookmark.Status,
			KeepReading: keepReading,
			MangaInfo: dtos.MangaInfo{
				Name:          mangaModel.Name,
				TotalChapters: pathModel.TotalChapters,
				LastUpdate:    pathModel.LastUpdate.Time(),
			},
		}

		bookmarks = append(bookmarks, bookmarkDetail)
	}

	result := dtos.UserBookmars{
		Bookmarks: bookmarks,
	}

	if totalBookmarks > 0 {
		result.TotalBookmarks = &totalBookmarks
	}

	return result, nil
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
	update := bson.D{
		{"$set", bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		}},
	}

	// Merge `updateBson` into the `$set` operator
	update[0].Value = mergeMaps(update[0].Value.(bson.M), updateBson)

	filter = bson.M{"_id": existingBookmark.Id}

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
	//updatedBookmark.KeepReading = false
	/*mangaId, _ := primitive.ObjectIDFromHex(updatedBookmark.MangaId)

	filter = bson.M{"_id": mangaId}
	mangaModel, code, err := repository.FindManga(filter)
	if err == nil && code == constants.NoError {
		updatedBookmark.KeepReading = updatedBookmark.Chapter < mangaModel.TotalChapters
	}*/

	return updatedBookmark, nil
}

func mergeMaps(m1, m2 bson.M) bson.M {
	for key, value := range m2 {
		m1[key] = value
	}
	return m1
}

/*func CheckForMangaUpdates(bookmarkId string) (dtos.Bookmark, error) {
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

	filter = bson.M{"_id": existingBookmark.PathId}

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
}*/

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

func findExistingBookmarkV2(pathID, userID primitive.ObjectID) (*models.Bookmark, error) {
	filter := bson.M{
		"pathId": pathID,
		"userId": userID,
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

	//TODO refactor

	/*mangaId, _ := primitive.ObjectIDFromHex(bookmark.MangaId)
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

	keepReading = bookmark.Chapter < mangaModel.TotalChapters*/

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
