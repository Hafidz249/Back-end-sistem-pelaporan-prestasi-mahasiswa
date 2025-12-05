package repository

import (
	"POJECT_UAS/model"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	DB        *sql.DB
	JWTSecret string
}

func (r *AuthRepository) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	var user model.Users

	// Query untuk mendukung login dengan username atau email
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active 
		FROM users 
		WHERE username = $1 OR email = $1
	`

	err := r.DB.QueryRow(query, req.Credential).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("kredensial salah")
		}
		return nil, err
	}

	// Cek status aktif user
	if !user.IsActive {
		return nil, errors.New("akun anda dinonaktifkan")
	}

	// Validasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("kredensial salah")
	}

	// Ambil role dan permissions
	profile, err := r.getUserProfile(user)
	if err != nil {
		return nil, err
	}

	// Generate JWT token dengan role dan permissions
	token, err := r.generateJWT(user, profile.Permissions)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:   token,
		Profile: *profile,
	}, nil
}

func (r *AuthRepository) getUserProfile(user model.Users) (*model.UserProfile, error) {
	// Ambil role info
	var role model.RoleInfo
	roleQuery := `
		SELECT id, name, description 
		FROM roles 
		WHERE id = $1
	`
	err := r.DB.QueryRow(roleQuery, user.RoleID).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, err
	}

	// Ambil permissions berdasarkan role
	permissionsQuery := `
		SELECT p.id, p.name, p.resource, p.action 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`
	rows, err := r.DB.Query(permissionsQuery, user.RoleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var perm model.Permission
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return &model.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		Role:        role,
		Permissions: permissions,
	}, nil
}

func (r *AuthRepository) generateJWT(user model.Users, permissions []model.Permission) (string, error) {
	// Convert permissions ke format untuk JWT
	permList := make([]map[string]string, len(permissions))
	for i, p := range permissions {
		permList[i] = map[string]string{
			"name":     p.Name,
			"resource": p.Resource,
			"action":   p.Action,
		}
	}

	claims := jwt.MapClaims{
		"user_id":     user.ID.String(),
		"username":    user.Username,
		"email":       user.Email,
		"role_id":     user.RoleID.String(),
		"permissions": permList,
		"exp":         time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(r.JWTSecret))
}