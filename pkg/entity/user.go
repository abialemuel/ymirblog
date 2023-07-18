package entity

// User is the entity that represents a user.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"email,required"`
}
