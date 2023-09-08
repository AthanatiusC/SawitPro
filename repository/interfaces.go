// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateUser(ctx context.Context, input CreateUserInput) (output User, err error)
	UpdateUserById(ctx context.Context, input UpdateUserInput) (output User, err error)
	GetUserById(ctx context.Context, id int) (output User, err error)
	GetUserByPhoneNumber(ctx context.Context, phone string) (output User, err error)
}
