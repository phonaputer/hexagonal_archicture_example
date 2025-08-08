package http_presentation

import (
	"encoding/json"
	"errors"
	"examplemodule/internal/exampleapp/logic"
	"log/slog"
	"net/http"
)

type JSONObjectHandler struct {
	crudService logic.JSONObjectCRUDService
	validator   JSONObjectRequestValidator
}

func NewJSONObjectHandler(
	crudService logic.JSONObjectCRUDService,
	validator JSONObjectRequestValidator,
) *JSONObjectHandler {
	return &JSONObjectHandler{
		crudService: crudService,
		validator:   validator,
	}
}

func (h *JSONObjectHandler) CreateNewObject(w http.ResponseWriter, r *http.Request) {

	// I/O. Receive input in HTTP format. Parse and validate it matches the expected HTTP request data model.

	requestData, err := h.validator.jsonObjectCreateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Map HTTP data model to business logic data model

	input := &logic.NewJSONObject{
		Object:     requestData.JSONObject,
		SFObjectID: requestData.SFObjectID,
		SchemaID:   requestData.SchemaID,
		UserID:     requestData.EndUserID,
	}

	// Call into business logic

	result, err := h.crudService.Create(r.Context(), input)
	if errors.Is(err, logic.ErrObjectAlreadyExists) { // Map logic error to HTTP error
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map business logic result data model to HTTP response data model

	responseBody := &jsonObjectResponse{
		ID:         result.ID,
		SchemaID:   result.SchemaID,
		SFObjectID: result.SFObjectID,
		EndUserID:  result.UserID,
		JSONObject: result.Object,
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

func (h *JSONObjectHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {

	// I/O. Receive input in HTTP format. Parse and validate it matches the expected HTTP request data model.

	requestData, err := h.validator.jsonObjectDeleteRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Don't need to map HTTP data model to business logic data model, because HTTP model already follows
	// business model (id is an integer).

	// Call into business logic

	err = h.crudService.Delete(r.Context(), requestData.ID)
	if errors.Is(err, logic.ErrObjectNotFound) { // Map logic error to HTTP error
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

func (h *JSONObjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	// I/O. Receive input in HTTP format. Parse and validate it matches the expected HTTP request data model.

	requestData, err := h.validator.jsonObjectGetByIDRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Don't need to map HTTP data model to business logic data model, because HTTP model already follows
	// business model (id is an integer).

	// Call into business logic

	result, err := h.crudService.GetByID(r.Context(), requestData.ID)
	if errors.Is(err, logic.ErrObjectNotFound) { // Map logic error to HTTP error
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map business logic result data model to HTTP response data model

	responseBody := &jsonObjectResponse{
		ID:         result.ID,
		SchemaID:   result.SchemaID,
		SFObjectID: result.SFObjectID,
		EndUserID:  result.UserID,
		JSONObject: result.Object,
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
