package services

import (
	"errors"
	"fmt"
	"manga-bookmarker-backend/dtos"
	"strings"
)

func CreateBookmark(data dtos.CreateBookmark) error {
	fmt.Println(fmt.Sprintf("%+v", data))
	//TODO obtain id of manga in URL

	// Find the index of "manga-"
	prefix := "manga-"
	idx := strings.Index(data.Url, prefix)
	if idx == -1 {
		fmt.Println("Prefix not found: ", data.Url)
		return errors.New("No se encontro el prefijo de manganato en la url")
	}

	// Extract the substring after "manga-"
	mangaID := data.Url[idx+len(prefix):]

	// Check if there's a "/" in the extracted ID and handle it accordingly
	if slashIdx := strings.Index(mangaID, "/"); slashIdx != -1 {
		mangaID = mangaID[:slashIdx] // Exclude the part after "/"
	}

	fmt.Println("Manga ID:", mangaID)

	//TODO check if id of manga exists

	//TODO create manga if not exists and obtain ObjectId

	//TODO check if the bookmark for that manga does not exists

	//TODO create the bookmark

	return nil
}
