package repository

import (
	"database/sql"
	"errors"

	"POJECT_UAS/model"

	"github.com/google/uuid"
)

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) Create(u *model.Users) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    query := `INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
    _, err := r.DB.Exec(query, u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID, u.ISActive, u.Created_at, u.Updated_at)
    return err
}

func (r *UserRepository) GetByEmail(email string) (*model.Users, error) {
    row := r.DB.QueryRow(`SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at FROM users WHERE email=$1`, email)
    var u model.Users
    err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.ISActive, &u.Created_at, &u.Updated_at)
    if err == sql.ErrNoRows {
        return nil, errors.New("not found")
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}
