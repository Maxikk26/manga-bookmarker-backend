package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bookmark struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MangaId  primitive.ObjectID `bson:"manga_id"`
	UserId   primitive.ObjectID `bson:"user_id"`
	Chapter  string             `bson:"chapter"`
	LastRead primitive.DateTime `bson:"last_read"`
}
