package exampleapp

import (
	"cloud.google.com/go/pubsub"
	"context"
	"examplemodule/internal/exampleapp/io/http_presentation"
	"examplemodule/internal/exampleapp/io/mysql_storage"
	"examplemodule/internal/exampleapp/io/pubsub_publisher"
	"examplemodule/internal/exampleapp/logic"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func RunApplication() {
	// Dependency injection

	db := sqlx.MustOpen("mysql", "EXAMPLE MYSQL ADDRESS")
	userStorage := mysql_storage.NewUserStorage(db)

	ps, err := pubsub.NewClient(context.Background(), "EXAMPLE GCP PROJECT ID")
	if err != nil {
		panic(err)
	}
	topic := ps.Topic("users")

	userPublisher := pubsub_publisher.NewUserEventPublisher(topic)

	service := logic.NewUserServiceLogic(userStorage, userPublisher)

	requestValidator := http_presentation.NewUserRequestValidatorImpl()
	userHandlers := http_presentation.NewUserHandler(service, requestValidator)

	// Configure and run simple HTTP server

	http.HandleFunc("POST /users", userHandlers.Create)
	http.HandleFunc("DELETE /users", userHandlers.DeleteByID)
	http.HandleFunc("GET /users", userHandlers.GetByID)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
