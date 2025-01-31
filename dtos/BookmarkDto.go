package dtos

import (
	"time"
)

type CreateBookmark struct {
	Url     string `json:"url"`
	Path    string `json:"path"`
	UserId  string `json:"userId,omitempty"`
	Chapter string `json:"chapter,omitempty"`
	Status  int    `json:"status,omitempty"`
	SiteId  string `json:"siteId,omitempty"`
}

type BookmarkMangaUpdate struct {
	Update      bool       `json:"update"`
	LastChapter string     `json:"lastChapter,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
}

type Bookmark struct {
	Id          string               `json:"id,omitempty"`
	MangaId     string               `json:"mangaId,omitempty"`
	Chapter     string               `json:"chapter,omitempty"`
	LastRead    time.Time            `json:"lastRead,omitempty"`
	Status      int                  `json:"status,omitempty"`
	MangaUpdate *BookmarkMangaUpdate `json:"mangaUpdate,omitempty"`
}

type BookmarkDetail struct {
	Id          string    `json:"id,omitempty"`
	Chapter     string    `json:"chapter,omitempty"`
	LastRead    time.Time `json:"lastRead,omitempty"`
	Status      int       `json:"status,omitempty"`
	KeepReading bool      `json:"keepReading,omitempty"`
	MangaInfo   MangaInfo `json:"mangaInfo,omitempty"`
}

type UserBookmars struct {
	TotalBookmarks *int64           `json:"totalBookmarks,omitempty"`
	Bookmarks      []BookmarkDetail `json:"bookmarks,omitempty"`
}

type BookmarkUpdate struct {
	Chapter  string    `json:"chapter,omitempty" bson:"chapter,omitempty"`
	LastRead time.Time `json:"lastRead,omitempty" bson:"lastRead,omitempty"`
	Status   int       `json:"status,omitempty" bson:"status,omitempty"`
}
