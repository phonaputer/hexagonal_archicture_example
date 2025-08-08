package pubsub_publisher

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"examplemodule/internal/exampleapp/logic"
	"fmt"
	"time"
)

type JSONObjectEventPublisher struct {
	topic *pubsub.Topic
}

func NewJSONObjectEventPublisher(topic *pubsub.Topic) *JSONObjectEventPublisher {
	return &JSONObjectEventPublisher{
		topic: topic,
	}
}

func (p *JSONObjectEventPublisher) PublishCreate(ctx context.Context, jsonObject *logic.JSONObject) error {
	// Map from business logic data model to GCP PubSub data model

	event := &createEventJSON{
		ID:           jsonObject.ID,
		SFObjectID:   jsonObject.SFObjectID,
		EndUserID:    jsonObject.UserID,
		CreationTime: time.Now().Format(time.RFC3339),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	// I/O. Publish event to GCP PubSub.

	_, err = p.topic.Publish(ctx, &pubsub.Message{Data: eventBytes}).Get(ctx)
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	// Nothing to do here since there is no result

	return nil
}

func (p *JSONObjectEventPublisher) PublishDelete(ctx context.Context, id int) error {
	// Map from business logic data model to GCP PubSub data model

	event := &deleteEventJSON{
		ID:           id,
		DeletionTime: time.Now().Format(time.RFC3339),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	// I/O. Publish event to GCP PubSub.

	_, err = p.topic.Publish(ctx, &pubsub.Message{Data: eventBytes}).Get(ctx)
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	// Nothing to do here since there is no result

	return nil
}
