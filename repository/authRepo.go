package repository

import (
	"POJECT_UAS/model"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB        *sql.DB
	JWTSecret string
}

func (r *UserRepository) Login(req model.Login) (string, error) {
	var user model.Users

	query := `
		SELECT id, username, password_hash, role_id, is_active 
		FROM users 
		WHERE username = $1
	`

	err := r.DB.QueryRow(query, req.Username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.RoleID,
		&user.ISActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("username salah")
		}
		return "", err
	}

	if !user.ISActive {
		return "", errors.New("akun anda dinonaktifkan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash))
	if err != nil {
		return "", errors.New("password salah")
	}

	token, err := r.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *UserRepository) generateJWT(user model.Users) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role_id":  user.RoleID.String(),
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(r.JWTSecret))
}