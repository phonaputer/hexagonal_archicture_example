package logic

import (
	"context"
	"errors"
	"fmt"
)

var ErrObjectNotFound = errors.New("object not found")
var ErrObjectAlreadyExists = errors.New("object already exists")

type JSONObjectCRUDService interface {
	Create(ctx context.Context, jsonObject *NewJSONObject) (*JSONObject, error)
	GetByID(ctx context.Context, id int) (*JSONObject, error)
	Delete(ctx context.Context, id int) error
}

type JSONObjectServiceLogic struct {
	storage   JSONObjectStorage
	publisher JSONObjectEventPublisher
}

func NewJSONObjectServiceLogic(
	storage JSONObjectStorage,
	publisher JSONObjectEventPublisher,
) *JSONObjectServiceLogic {
	return &JSONObjectServiceLogic{
		storage:   storage,
		publisher: publisher,
	}
}

func (s *JSONObjectServiceLogic) Create(ctx context.Context, jsonObject *NewJSONObject) (*JSONObject, error) {
	exists, err := s.storage.ExistsBySFObjectID(ctx, jsonObject.SFObjectID)
	if err != nil {
		return nil, fmt.Errorf("exists by sf object id: %w", err)
	}
	if exists {
		return nil, ErrObjectAlreadyExists
	}

	newID, err := s.storage.Create(ctx, jsonObject)
	if err != nil {
		return nil, fmt.Errorf("create object: %w", err)
	}

	result := &JSONObject{
		ID:         newID,
		Object:     jsonObject.Object,
		SFObjectID: jsonObject.SFObjectID,
		SchemaID:   jsonObject.SchemaID,
		UserID:     jsonObject.UserID,
	}

	err = s.publisher.PublishCreate(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("publish create object: %w", err)
	}

	return result, nil
}

func (s *JSONObjectServiceLogic) GetByID(ctx context.Context, id int) (*JSONObject, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *JSONObjectServiceLogic) Delete(ctx context.Context, id int) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	err = s.publisher.PublishDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("publish delete object: %w", err)
	}

	return nil
}
