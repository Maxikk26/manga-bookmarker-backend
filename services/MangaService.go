package services

import (
	"errors"
	"fmt"
	"manga-bookmarker-backend/constants"
	"manga-bookmarker-backend/dtos"
	"manga-bookmarker-backend/repository"
	"manga-bookmarker-backend/utils"
)

func CreateManga() {

}

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
