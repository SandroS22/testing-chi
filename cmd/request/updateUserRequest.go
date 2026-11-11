package request

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

func NewUpdateUserRequest(name, email, password, role string) *UpdateUserRequest {
	return &UpdateUserRequest{
		Name:     name,
		Email:    email,
		User:     email,
		Password: password,
		Role:     role,
	}
}
