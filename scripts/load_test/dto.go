package main

type CreateUserRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type UserResponse struct {
	ID int `json:"id"`
}

type CreateTaskRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	AuthorUserID int    `json:"author_user_id"`
}

type TaskResponse struct {
	ID      int `json:"id"`
	Version int `json:"version"`
}

type PatchTaskRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
