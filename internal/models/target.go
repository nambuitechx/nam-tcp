package models

type TargetModel struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

type CreateTargetPayload struct {	
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
}
