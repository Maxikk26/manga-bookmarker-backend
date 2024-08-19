package services

import (
	"fmt"
	"manga-bookmarker-backend/repository"
	"os"
)

func ScrappingJob() {
	mangas, _, err := repository.AllMangas()
	if err != nil {
		fmt.Println("Error getting all mangas")
	}

	for _, manga := range mangas {
		url := os.Getenv("MANGANATO_URL") + manga.Identifier
		go AsyncUpdatesScrapping(url, manga)
	}

}
