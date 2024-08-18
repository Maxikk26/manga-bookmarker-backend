package constants

const (
	//Mongo error handling
	NoError         = -1
	Other           = 0
	NoDocumentFound = 1

	//TODO endpoints que devuelva estatus de lectura

	//Bookmark status
	Reading        = 1
	PlanningToRead = 2
	OnHold         = 3
	Dropped        = 4
)
