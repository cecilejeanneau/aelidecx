package main

type Note struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Folder    string `json:"folder"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type noteInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Folder  string `json:"folder"`
}
