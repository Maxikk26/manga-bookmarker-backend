package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manga struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Identifier    string             `bson:"identifier"`
	Name          string             `bson:"name"`
	Cover         string             `bson:"cover"`
	TotalChapters string             `bson:"totalChapters"`
	LastUpdate    primitive.DateTime `bson:"lastUpdate"`
}
