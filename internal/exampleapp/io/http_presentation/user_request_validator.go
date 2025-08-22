package http_presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UserRequestValidator interface {
	userCreateRequest(r *http.Request) (*userCreateRequest, error)
	userDeleteRequest(r *http.Request) (*userDeleteRequest, error)
	userGetByIDRequest(r *http.Request) (*userGetByIDRequest, error)
}

type UserRequestValidatorImpl struct{}

func NewUserRequestValidatorImpl() *UserRequestValidatorImpl {
	return &UserRequestValidatorImpl{}
}

func (v *UserRequestValidatorImpl) userCreateRequest(r *http.Request) (*userCreateRequest, error) {
	var reqData struct {
		FirstName    *string `json:"firstName"`
		LastName     *string `json:"lastName"`
		EmailAddress *string `json:"emailAddress"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		return nil, errors.New("request body is not valid json")
	}

	err = verifyNotNil(err, reqData.FirstName, "firstName is required")
	err = verifyNotNil(err, reqData.LastName, "lastName is required")
	err = verifyNotNil(err, reqData.EmailAddress, "emailAddress is required")

	if err != nil {
		return nil, err
	}

	return &userCreateRequest{
		FirstName:    *reqData.FirstName,
		LastName:     *reqData.LastName,
		EmailAddress: *reqData.EmailAddress,
	}, nil
}

func (v *UserRequestValidatorImpl) userDeleteRequest(r *http.Request) (*userDeleteRequest, error) {
	var reqData struct {
		ID *int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		return nil, errors.New("request body is not valid json")
	}

	err = verifyNotNil(err, reqData.ID, "id is required")

	if err != nil {
		return nil, err
	}

	return &userDeleteRequest{
		ID: *reqData.ID,
	}, nil
}

func (v *UserRequestValidatorImpl) userGetByIDRequest(r *http.Request) (*userGetByIDRequest, error) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		return nil, errors.New("id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("id is not a valid integer")
	}

	return &userGetByIDRequest{
		ID: id,
	}, nil
}

func verifyNotNil[T any](err error, value *T, msg string) error {
	if err != nil {
		return err
	}

	if value == nil {
		return errors.New(msg)
	}

	return nil
}
