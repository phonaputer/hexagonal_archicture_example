package logic

import (
	"context"
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")

type UserService interface {
	Create(ctx context.Context, user *NewUser) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Delete(ctx context.Context, id int) error
}

type UserServiceLogic struct {
	storage   UserStorage
	publisher UserEventPublisher
}

func NewUserServiceLogic(
	storage UserStorage,
	publisher UserEventPublisher,
) *UserServiceLogic {
	return &UserServiceLogic{
		storage:   storage,
		publisher: publisher,
	}
}

func (s *UserServiceLogic) Create(ctx context.Context, newUser *NewUser) (*User, error) {
	exists, err := s.storage.ExistsByEmailAddress(ctx, newUser.EmailAddress)
	if err != nil {
		return nil, fmt.Errorf("exists by email address: %w", err)
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	newID, err := s.storage.Create(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	result := &User{
		ID:           newID,
		EmailAddress: newUser.EmailAddress,
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
	}

	err = s.publisher.PublishCreate(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("publish create: %w", err)
	}

	return result, nil
}

func (s *UserServiceLogic) GetByID(ctx context.Context, id int) (*User, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *UserServiceLogic) Delete(ctx context.Context, id int) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	err = s.publisher.PublishDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("publish delete: %w", err)
	}

	return nil
}
