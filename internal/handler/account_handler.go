package handler

import (
	"banking-api/internal/middleware"
	"banking-api/internal/models"
	"banking-api/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	AccountService *service.AccountService
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{AccountService: service}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ошибка токена", http.StatusUnauthorized)
		return
	}

	account, err := h.AccountService.CreateAccount(userID)
	if err != nil {
		http.Error(w, "не удалось создать счёт", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"id":        account.ID,
		"balance":   account.Balance,
		"createdAt": account.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AccountHandler) TopUp(w http.ResponseWriter, r *http.Request) {
	var req models.TopUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "невалидный JSON", http.StatusBadRequest)
		return
	}

	userIDStr := middleware.GetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ошибка токена", http.StatusUnauthorized)
		return
	}

	if err := h.AccountService.TopUp(userID, req.AccountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "невалидный JSON", http.StatusBadRequest)
		return
	}

	userIDStr := middleware.GetUserID(r.Context())
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ошибка токена", http.StatusUnauthorized)
		return
	}

	err = h.AccountService.TransferFunds(userID, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *AccountHandler) TransferByUsernames(w http.ResponseWriter, r *http.Request) {
	var req models.TransferByUsernamesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "невалидный JSON", http.StatusBadRequest)
		return
	}

	err := h.AccountService.TransferBetweenUsers(req.FromUsername, req.ToUsername, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
