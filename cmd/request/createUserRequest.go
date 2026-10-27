package request

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func NewCreateUserRequest(name, email, password, role string) *CreateUserRequest {
	return &CreateUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
}
