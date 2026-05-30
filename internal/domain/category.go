package domain

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CategoryPayload struct {
	Name  string  `json:"name"  binding:"required"`
	Color *string `json:"color"`
}

type CategoryRepository interface {
	FindAll() ([]Category, error)
	FindByID(id int) (*Category, error)
	Create(payload CategoryPayload) (*Category, error)
	Delete(id int) error
}
