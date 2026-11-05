package handler

import (
	"bank/internal/calculator"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Handler - обработчик калькулятора. Хранит сервис, у которого просит вычисления.
type Handler struct {
	service *calculator.Service
}

// NewHandler - конструктор: принимает сервис и возвращает готовый обработчик.
func NewHandler(service *calculator.Service) *Handler {
	return &Handler{service: service}
}

// Calculate - обработчик POST /calculate: читает JSON, считает, отвечает JSON.
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	// Калькулятор принимает данные, поэтому разрешаем только POST.
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Превращаем JSON из тела запроса в структуру Request.
	var req calculator.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	// Вся работа - в сервисе, обработчик только принимает и отвечает.
	resp, err := h.service.Calculate(req)
	if err != nil {
		// Эти три ошибки - вина клиента (плохие данные), отвечаем 400.
		if errors.Is(err, calculator.ErrUnknowOperation) ||
			errors.Is(err, calculator.ErrDivzionByZero) ||
			errors.Is(err, calculator.ErrNegativeAmount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Всё остальное - сбой у нас: подробности в лог, клиенту коротко
		log.Printf("error calculating: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	// Говорим клиенту, что в ответе JSON, и отправляем результат.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
