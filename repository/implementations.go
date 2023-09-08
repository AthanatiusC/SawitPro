package repository

import (
	"context"
)

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (output User, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	query := `INSERT INTO users(name, phone, password) values($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, query, input.Name, input.Phone, input.Password).Scan(&output.Id)
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

func (r *Repository) UpdateUserById(ctx context.Context, input UpdateUserInput) (output User, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return
	}
	defer tx.Commit()

	query := `UPDATE users SET name=$1, phone=$2, updated_at=NOW() WHERE id = $3 RETURNING name,phone`
	err = r.Db.QueryRowContext(ctx, query, input.Name, input.Phone, input.Id).Scan(
		&output.Name,
		&output.Phone,
	)
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

func (r *Repository) GetUserById(ctx context.Context, id int) (output User, err error) {
	query := `SELECT u.id, u.name, u.phone, u.password, u.updated_at, u.created_at FROM users u WHERE u.id = $1`
	err = r.Db.QueryRowContext(ctx, query, id).Scan(
		&output.Id,
		&output.Name,
		&output.Phone,
		&output.Password,
		&output.UpdatedAt,
		&output.CreatedAt,
	)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phone string) (output User, err error) {
	query := `SELECT u.id, u.name, u.phone, u.password, u.updated_at, u.created_at FROM users u WHERE u.phone = $1`
	err = r.Db.QueryRowContext(ctx, query, phone).Scan(
		&output.Id,
		&output.Name,
		&output.Phone,
		&output.Password,
		&output.UpdatedAt,
		&output.CreatedAt,
	)
	if err != nil {
		return
	}
	return
}
