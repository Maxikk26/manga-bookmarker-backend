package dtos

import (
	"time"
)

type MangaScrapperData struct {
	Name          string    `json:"name"`
	Cover         string    `json:"cover"`
	TotalChapters string    `json:"totalChapters"`
	LastUpdate    time.Time `bson:"lastUpdate"`
}
