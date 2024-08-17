package dtos

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type Manga struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`        // MongoDB ObjectID
	Title       string             `bson:"title"`                // Title of the manga
	Author      string             `bson:"author"`               // Author of the manga
	Identifier  string             `bson:"identifier"`           // Unique identifier for the manga
	Description string             `bson:"description"`          // Description of the manga
	Genre       []string           `bson:"genre"`                // List of genres
	CoverURL    string             `bson:"cover_url"`            // URL to the cover image
	Chapters    int                `bson:"chapters"`             // Number of chapters available
	Status      string             `bson:"status"`               // Status of the manga (e.g., ongoing, completed)
	CreatedAt   primitive.DateTime `bson:"created_at,omitempty"` // Date and time the manga was added
	UpdatedAt   primitive.DateTime `bson:"updated_at,omitempty"` // Date and time the manga was last updated
}
