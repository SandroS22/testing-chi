package response

type CreateUserResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func NewCreateUserResponse(id int, name, email, password, role string) *CreateUserResponse {
	return &CreateUserResponse{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
}
