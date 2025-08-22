package http_presentation

import (
	"encoding/json"
	"errors"
	"examplemodule/internal/exampleapp/logic"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	crudService logic.UserService
	validator   UserRequestValidator
}

func NewUserHandler(
	crudService logic.UserService,
	validator UserRequestValidator,
) *UserHandler {
	return &UserHandler{
		crudService: crudService,
		validator:   validator,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	// I/O. Receive input in HTTP format. Parse and validate it matches the expected HTTP request data model.

	requestData, err := h.validator.userCreateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Map HTTP data model to business logic data model

	input := &logic.NewUser{
		EmailAddress: requestData.EmailAddress,
		FirstName:    requestData.FirstName,
		LastName:     requestData.LastName,
	}

	// Call into business logic

	result, err := h.crudService.Create(r.Context(), input)
	if errors.Is(err, logic.ErrUserAlreadyExists) { // Map logic error to HTTP error
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map business logic result data model to HTTP response data model

	responseBody := &userResponse{
		ID:           result.ID,
		FirstName:    result.FirstName,
		LastName:     result.LastName,
		EmailAddress: result.EmailAddress,
	}

	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// I/O. Write HTTP response data model to HTTP request sender

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBytes)
	if err != nil {
		slog.Error("Error writing response", slog.String("err", err.Error()))
	}
}

func (h *UserHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {

	// I/O. Receive input in HTTP format. Parse and validate it matches the expected HTTP request data model.

	requestData, err := h.validator.userDeleteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Don't need to map HTTP data model to business logic data model, because HTTP model already follows
	// business model (id is an integer).

	// Call into business logic

	err = h.crudService.Delete(r.Context(), requestData.ID)
	if errors.Is(err, logic.ErrUserNotFound) { // Map logic error to HTTP error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No result to map

	// I/O. Write HTTP response to HTTP request sender

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	// I/O. Receive input in HTTP format. Parse and validate it matches the expected HTTP request data model.

	requestData, err := h.validator.userGetByIDRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Don't need to map HTTP data model to business logic data model, because HTTP model already follows
	// business model (id is an integer).

	// Call into business logic

	result, err := h.crudService.GetByID(r.Context(), requestData.ID)
	if errors.Is(err, logic.ErrUserNotFound) { // Map logic error to HTTP error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map business logic result data model to HTTP response data model

	responseBody := &userResponse{
		ID:           result.ID,
		FirstName:    result.FirstName,
		LastName:     result.LastName,
		EmailAddress: result.EmailAddress,
	}

	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// I/O. Write HTTP response data model to HTTP request sender

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		slog.Error("Error writing response", slog.String("err", err.Error()))
	}
}
