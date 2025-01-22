package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Path struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	SiteId        primitive.ObjectID `bson:"siteId" json:"siteId"`
	MangaId       primitive.ObjectID `bson:"mangaId" json:"mangaId"`
	Path          string             `bson:"path" json:"path"`
	TotalChapters string             `bson:"totalChapters"`
	LastUpdate    primitive.DateTime `bson:"lastUpdate"`
}
