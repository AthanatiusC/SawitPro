// This file contains types that are used in the repository layer.
package repository

import "time"

type CreateUserInput struct {
	Name     string
	Phone    string
	Password string
}

type UpdateUserInput struct {
	Id    int
	Name  string
	Phone string
}

type User struct {
	Id        int
	Name      string
	Phone     string
	Password  string
	UpdatedAt time.Time
	CreatedAt time.Time
}
