package services

import (
	"errors"
	"fmt"
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

	fmt.Println("manga", manga)

	if errorType != constants.NoDocumentFound {
		fmt.Println("Error type NoDocumentFound after checking for err variable:", err.Error())
		return errors.New("Ocurrio un error buscando el manga")
	}

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

	fmt.Println(fmt.Sprintf("%+v", manga))

	id, err := repository.CreateManga(manga)
	if err != nil {
		fmt.Println("Error creating manga: ", err.Error())
		return errors.New("Ocurrio un error creando el manga")
	}

	fmt.Println(id)

	//TODO check if the bookmark for that manga does not exists

	//TODO create the bookmark

	return nil
}
