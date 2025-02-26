package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO manga status (ongoing,completed...)
// TODO manga genres

type Manga struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name"`
	Cover     string             `bson:"cover"`
	UpdatedAt primitive.DateTime `bson:"updatedAt"`
}
