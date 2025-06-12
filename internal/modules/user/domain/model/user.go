package model

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Password  string `json:"-"` // El guión hace que no se serialice en JSON
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
