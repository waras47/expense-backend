package dto

type CategoryPayload struct {
	Name  string  `json:"name"  binding:"required"`
	Color *string `json:"color"`
}
