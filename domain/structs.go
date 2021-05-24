package domain

//easyjson:json
type User struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Browsers []string `json:"browsers"`
}
