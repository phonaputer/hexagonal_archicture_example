package logic

import "context"

type UserStorage interface {
	Create(ctx context.Context, user *NewUser) (int, error)
	ExistsByEmailAddress(ctx context.Context, emailAddress string) (bool, error)
	GetByID(ctx context.Context, id int) (*User, error)
	Delete(ctx context.Context, id int) error
}
