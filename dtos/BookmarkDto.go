package dtos

import (
	"time"
)

type CreateBookmark struct {
	Url     string `json:"url"`
	UserId  string `json:"userId,omitempty"`
	Chapter string `json:"chapter,omitempty"`
	Status  int    `json:"status,omitempty"`
}

type Bookmark struct {
	Id       string    `json:"id,omitempty"`
	Chapter  string    `json:"chapter,omitempty"`
	LastRead time.Time `json:"lastRead,omitempty"`
	Status   int       `json:"status,omitempty"`
}
