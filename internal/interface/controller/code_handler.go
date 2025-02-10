package controller

import (
	"compiler-playground-api/internal/entity"
	"compiler-playground-api/internal/usecase"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CodeHandler struct {
	useCase usecase.CodeUseCase
}

func NewCodeHandler(uc *usecase.CodeUseCase) *CodeHandler {
	return &CodeHandler{useCase: *uc}
}

func (h *CodeHandler) SaveCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var code entity.Code
	if err := json.NewDecoder(r.Body).Decode(&code); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	code.ID = uuid.New().String()
	code.CreatedAt = time.Now()

	id, err := h.useCase.SaveCode(&code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Code saved successfully",
		"id":         id,
		"created_at": code.CreatedAt,
	})
}

func (h *CodeHandler) ExecuteCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	output, err := h.useCase.ExecuteCode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Code executed successfully",
		"output":  output,
	})
}
