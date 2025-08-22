package logic

import "context"

type UserEventPublisher interface {
	PublishCreate(ctx context.Context, user *User) error
	PublishDelete(ctx context.Context, id int) error
}
