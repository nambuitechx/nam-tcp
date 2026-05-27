package models

type UserModel struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

type CreateUserPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}
