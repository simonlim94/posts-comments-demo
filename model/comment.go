package model

type Comment struct {
	ID     uint32 `json:"id"`
	PostID uint32 `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}
