package logic

import "context"

type JSONObjectEventPublisher interface {
	PublishCreate(ctx context.Context, jsonObject *JSONObject) error
	PublishDelete(ctx context.Context, id int) error
}
