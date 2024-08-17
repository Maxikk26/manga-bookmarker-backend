package dtos

import (
	"time"
)

type MangaScrapperData struct {
	Name          string    `json:"name"`
	Cover         string    `json:"cover"`
	TotalChapters string    `json:"totalChapters"`
	LastUpdate    time.Time `json:"lastUpdate"`
}

type MangaInfo struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	TotalChapters string    `json:"totalChapters"`
	LastUpdate    time.Time `json:"lastUpdate"`
}
