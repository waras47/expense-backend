package dto

type CreateCategoryPayload struct {
	Name  string  `json:"name"  binding:"required,min=1,max=100"`
	Color *string `json:"color" binding:"omitempty,hexcolor"`
}

type UpdateCategoryPayload struct {
	Name  *string `json:"name"  binding:"omitempty,min=1,max=100"`
	Color *string `json:"color" binding:"omitempty,hexcolor"`
}
