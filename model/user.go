package model

import "time"

type Users struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	FindByID(id int) (*Users, error)
	FindByEmail(email string) (*Users, error)
	FindAll() ([]Users, error)
	Create(user *Users) error
	Update(user *Users) error
	Delete(id int) error
	Count(search string) (int, error)
}