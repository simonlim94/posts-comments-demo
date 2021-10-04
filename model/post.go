package model

type Post struct {
	ID     uint32 `json:"id"`
	UserID uint32 `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
