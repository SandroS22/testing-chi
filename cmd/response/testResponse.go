package response

type TestResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	User string `json:"user"`
	Role string `json:"role"`
}

func NewTestResponse(id int, name, user, role string) *TestResponse {
	return &TestResponse{
		Id:   id,
		Name: name,
		User: user,
		Role: role,
	}
}
