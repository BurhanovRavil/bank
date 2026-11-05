package handler

import (
	"bank/internal/calculator"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// Handler - обработчик калькулятора. Для работы ему нужны сервис (он считает)
// и логгер (он пишет, что происходит).
type Handler struct {
	service *calculator.Service
	logger  *slog.Logger
}

// NewHandler получает сервис и логгер снаружи, из main.go. Так настройками логов
// управляет main.go, а обработчик о них ничего не знает.
func NewHandler(service *calculator.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Calculate - обработчик POST /calculate: читает JSON, считает, отвечает JSON.
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {

	// Превращаем JSON из тела запроса в структуру Request.
	var req calculator.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid json in request", "err", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Подробности запроса интересны только при отладке, поэтому Debug.
	h.logger.Debug("calculate request", "a", req.A, "b", req.B, "op", req.Op)

	// Вся работа - в сервисе, обработчик только принимает и отвечает.
	resp, err := h.service.Calculate(req)
	if err != nil {
		// Эти три ошибки - вина клиента (плохие данные), отвечаем 400.
		if errors.Is(err, calculator.ErrUnknownOperation) ||
			errors.Is(err, calculator.ErrDivisionByZero) ||
			errors.Is(err, calculator.ErrNegativeAmount) {
			h.logger.Warn("calculate rejested", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Всё остальное - сбой у нас: подробности в лог, клиенту коротко
		h.logger.Error("calculate failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.logger.Info("calculated", "op", req.Op, "result", resp.Result)

	// Говорим клиенту, что в ответе JSON, и отправляем результат.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("encode response failed", "err", err)
	}
}
