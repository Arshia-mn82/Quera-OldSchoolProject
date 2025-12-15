package dto

type CreatePersonDTO struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}

type WhoAmIDTO struct {
	ID uint `json:"id,omitempty"`
}
