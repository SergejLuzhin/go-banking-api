package service

import (
	"banking-api/internal/config"
	"banking-api/internal/models"
	"banking-api/internal/repository"
	"errors"
	"fmt"
)

type AccountService struct {
	Repo         *repository.AccountRepository
	UserRepo     *repository.UserRepository
	EmailService *EmailService
}

func NewAccountService(repo *repository.AccountRepository, userRepo *repository.UserRepository, email *EmailService) *AccountService {
	return &AccountService{
		Repo:         repo,
		UserRepo:     userRepo,
		EmailService: email,
	}
}

func (s *AccountService) CreateAccount(userID int64) (*models.Account, error) {
	account, err := s.Repo.CreateAccount(userID)
	if err != nil {
		config.Log.Errorf("Ошибка создания счёта: %v", err)
		return nil, err
	}
	config.Log.Infof("Создан счёт ID=%d для пользователя ID=%d", account.ID, userID)
	return account, nil
}

func (s *AccountService) TopUp(userID, accountID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть больше 0")
	}
	err := s.Repo.TopUpAccount(accountID, userID, amount)
	if err != nil {
		config.Log.Errorf("Ошибка пополнения счёта %d: %v", accountID, err)
		return err
	}
	config.Log.Infof("Счёт ID=%d пополнен на %.2f пользователем ID=%d", accountID, amount, userID)
	return nil
}

func (s *AccountService) TransferFunds(userID, fromID, toID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	if fromID == toID {
		return errors.New("нельзя переводить самому себе")
	}
	err := s.Repo.TransferFunds(fromID, toID, userID, amount)
	if err != nil {
		config.Log.Errorf("Ошибка перевода: from %d to %d amount %.2f: %v", fromID, toID, amount, err)
		return err
	}
	config.Log.Infof("Перевод %.2f со счёта %d на счёт %d (user ID=%d)", amount, fromID, toID, userID)

	// Уведомление по email (опционально)
	if s.EmailService != nil {
		// Получаем email получателя
		toUserID, err := s.Repo.GetUserIDByAccountID(toID)
		if err == nil {
			receiver, err := s.UserRepo.GetUserByID(toUserID)
			if err == nil {
				body := fmt.Sprintf("<h3>Вам поступил перевод на сумму %.2f RUB</h3>", amount)
				_ = s.EmailService.SendEmail(receiver.Email, "Вы получили перевод", body)
			}
		}
	}

	return nil
}

func (s *AccountService) TransferToUsername(fromUserID, fromAccountID int64, toUsername string, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	if toUsername == "" {
		return errors.New("получатель не указан")
	}

	toUserID, err := s.UserRepo.GetUserIDByUsername(toUsername)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	toAccountID, err := s.Repo.GetFirstAccountByUserID(toUserID)
	if err != nil {
		return errors.New("счёт получателя не найден")
	}

	return s.TransferFunds(fromUserID, fromAccountID, toAccountID, amount)
}

func (s *AccountService) TransferBetweenUsers(fromUsername, toUsername string, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма должна быть положительной")
	}
	if fromUsername == toUsername {
		return errors.New("нельзя переводить самому себе")
	}

	fromUserID, err := s.UserRepo.GetUserIDByUsername(fromUsername)
	if err != nil {
		return errors.New("отправитель не найден")
	}
	toUserID, err := s.UserRepo.GetUserIDByUsername(toUsername)
	if err != nil {
		return errors.New("получатель не найден")
	}

	fromAccountID, err := s.Repo.GetFirstAccountByUserID(fromUserID)
	if err != nil {
		return errors.New("счёт отправителя не найден")
	}
	toAccountID, err := s.Repo.GetFirstAccountByUserID(toUserID)
	if err != nil {
		return errors.New("счёт получателя не найден")
	}

	return s.TransferFunds(fromUserID, fromAccountID, toAccountID, amount)
}
