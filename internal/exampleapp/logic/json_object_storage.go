package logic

import "context"

type JSONObjectStorage interface {
	Create(ctx context.Context, jsonObject *NewJSONObject) (int, error)
	ExistsBySFObjectID(ctx context.Context, objectID string) (bool, error)
	GetByID(ctx context.Context, id int) (*JSONObject, error)
	Delete(ctx context.Context, id int) error
}
