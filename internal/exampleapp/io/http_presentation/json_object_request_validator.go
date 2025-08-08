package http_presentation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type JSONObjectRequestValidator interface {
	jsonObjectCreateRequest(r *http.Request) (*jsonObjectCreateRequest, error)
	jsonObjectDeleteRequest(r *http.Request) (*jsonObjectDeleteRequest, error)
	jsonObjectGetByIDRequest(r *http.Request) (*jsonObjectGetByIDRequest, error)
}

type JSONObjectRequestValidatorImpl struct{}

func NewJSONObjectRequestValidatorImpl() *JSONObjectRequestValidatorImpl {
	return &JSONObjectRequestValidatorImpl{}
}

func (v *JSONObjectRequestValidatorImpl) jsonObjectCreateRequest(r *http.Request) (*jsonObjectCreateRequest, error) {
	var reqData struct {
		SchemaID   *string `json:"schemaId"`
		SFObjectID *string `json:"sfObjectId"`
		EndUserID  *string `json:"endUserId"`
		JSONObject *string `json:"jsonObject"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		return nil, errors.New("request body is not valid json")
	}

	err = verifyNotNil(err, reqData.SFObjectID, "sfObjectId is required")
	err = verifyNotNil(err, reqData.SchemaID, "schemaId is required")
	err = verifyNotNil(err, reqData.EndUserID, "endUserId is required")
	err = verifyNotNil(err, reqData.JSONObject, "jsonObject is required")

	if err != nil {
		return nil, err
	}

	return &jsonObjectCreateRequest{
		SchemaID:   *reqData.SchemaID,
		SFObjectID: *reqData.SFObjectID,
		EndUserID:  *reqData.EndUserID,
		JSONObject: *reqData.JSONObject,
	}, nil
}

func (v *JSONObjectRequestValidatorImpl) jsonObjectDeleteRequest(r *http.Request) (*jsonObjectDeleteRequest, error) {
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

	return &jsonObjectDeleteRequest{
		ID: *reqData.ID,
	}, nil
}

func (v *JSONObjectRequestValidatorImpl) jsonObjectGetByIDRequest(r *http.Request) (*jsonObjectGetByIDRequest, error) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		return nil, errors.New("id is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("id is not a valid integer")
	}

	return &jsonObjectGetByIDRequest{
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
