package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//TODO atributo SourceID para parametrizar tags del html

type Bookmark struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MangaId   primitive.ObjectID `bson:"mangaId"`
	UserId    primitive.ObjectID `bson:"userId"`
	Chapter   string             `bson:"chapter"`
	LastRead  primitive.DateTime `bson:"lastRead"`
	Status    int                `bson:"status"`
	UpdatedAt primitive.DateTime `bson:"updatedAt"`
}
