package repository

import (
	"banking-api/internal/models"
	"database/sql"
	//"errors"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) IsEmailOrUsernameTaken(email, username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2`
	var count int
	err := r.DB.QueryRow(query, email, username).Scan(&count)
	return count > 0, err
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.DB.QueryRow(query, user.Email, user.Username, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, created_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserIDByUsername(username string) (int64, error) {
	var userID int64
	err := r.DB.QueryRow(`SELECT id FROM users WHERE username = $1`, username).Scan(&userID)
	return userID, err
}

func (r *UserRepository) GetUserByID(userID int64) (*models.User, error) {
	query := `SELECT id, email, username FROM users WHERE id = $1`
	row := r.DB.QueryRow(query, userID)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
