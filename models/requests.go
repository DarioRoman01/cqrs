package models

type SignupLoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
