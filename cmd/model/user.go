package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

func NewUser(id int, name, email, user, password, role string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		User:     user,
		Password: password,
		Role:     role,
	}
}
