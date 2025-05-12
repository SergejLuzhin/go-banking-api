package service

import (
	"banking-api/internal/config"
	"banking-api/internal/models"
	"banking-api/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) RegisterUser(req *models.RegisterRequest) (*models.User, error) {
	exists, err := s.UserRepo.IsEmailOrUsernameTaken(req.Email, req.Username)
	if err != nil {
		config.Log.Errorf("Ошибка проверки уникальности: %v", err)
		return nil, err
	}
	if exists {
		return nil, errors.New("email или username уже используется")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		config.Log.Errorf("Ошибка хеширования пароля: %v", err)
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashed),
	}

	if err := s.UserRepo.CreateUser(user); err != nil {
		config.Log.Errorf("Ошибка создания пользователя: %v", err)
		return nil, err
	}

	config.Log.Infof("Пользователь зарегистрирован: %s", req.Email)
	return user, nil
}

func (s *AuthService) LoginUser(req *models.LoginRequest, jwtSecret string) (string, error) {
	user, err := s.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		config.Log.Warnf("Ошибка входа (email не найден): %s", req.Email)
		return "", errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		config.Log.Warnf("Ошибка входа (неверный пароль): %s", req.Email)
		return "", errors.New("неверный email или пароль")
	}

	token, err := GenerateJWT(user.ID, jwtSecret)
	if err != nil {
		config.Log.Errorf("Ошибка генерации токена: %v", err)
		return "", err
	}

	config.Log.Infof("Пользователь вошёл: %s", user.Email)
	return token, nil
}
