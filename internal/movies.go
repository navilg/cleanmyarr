package internal

type MovieFileDetail struct {
	MovieId     int    `json:"movieId"`
	MovieFileId int    `json:"id"`
	Size        int    `json:"size"`
	DateAdded   string `json:"dateAdded"`
}

type Movie struct {
	Title     string          `json:"title"`
	HasFile   bool            `json:"hasFile"`
	Tags      []int           `json:"tags"`
	ID        int             `json:"id"`
	MovieFile MovieFileDetail `json:"movieFile"`
}

type Tag struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type MovieImportEvent struct {
	MovieId int    `json:"movieId"`
	Date    string `json:"date"`
}
