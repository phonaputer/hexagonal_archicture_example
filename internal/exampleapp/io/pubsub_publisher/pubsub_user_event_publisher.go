package pubsub_publisher

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"examplemodule/internal/exampleapp/logic"
	"fmt"
	"time"
)

type UserEventPublisher struct {
	topic *pubsub.Topic
}

func NewUserEventPublisher(topic *pubsub.Topic) *UserEventPublisher {
	return &UserEventPublisher{
		topic: topic,
	}
}

func (p *UserEventPublisher) PublishCreate(ctx context.Context, user *logic.User) error {
	// Map from business logic data model to GCP PubSub data model

	event := &createUserEventJSON{
		ID:           user.ID,
		EmailAddress: user.EmailAddress,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
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

func (p *UserEventPublisher) PublishDelete(ctx context.Context, id int) error {
	// Map from business logic data model to GCP PubSub data model

	event := &deleteUserEventJSON{
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
