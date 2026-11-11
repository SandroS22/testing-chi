package response

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	User  string `json:"user"`
	Role  string `json:"role"`
}

func NewUserResponse(id int, name, email, role string) *UserResponse {
	return &UserResponse{
		ID:    id,
		Name:  name,
		Email: email,
		User:  email,
		Role:  role,
	}
}
