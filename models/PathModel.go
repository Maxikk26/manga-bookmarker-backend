package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Path struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	SiteId        primitive.ObjectID `json:"siteId"`
	MangaId       primitive.ObjectID `json:"mangaId"`
	Path          string             `json:"path"`
	TotalChapters string             `bson:"totalChapters"`
	LastUpdate    primitive.DateTime `bson:"lastUpdate"`
}
