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
	jsonObjectStorage := mysql_storage.NewJSONObjectStorage(db)

	ps, err := pubsub.NewClient(context.Background(), "EXAMPLE GCP PROJECT ID")
	if err != nil {
		panic(err)
	}
	topic := ps.Topic("json_objects")

	jsonObjectPublisher := pubsub_publisher.NewJSONObjectEventPublisher(topic)

	service := logic.NewJSONObjectServiceLogic(jsonObjectStorage, jsonObjectPublisher)

	requestValidator := http_presentation.NewJSONObjectRequestValidatorImpl()
	handlers := http_presentation.NewJSONObjectHandler(service, requestValidator)

	// Configure and run simple HTTP server

	http.HandleFunc("POST /json_objects", handlers.CreateNewObject)
	http.HandleFunc("DELETE /json_objects", handlers.DeleteByID)
	http.HandleFunc("GET /json_objects", handlers.GetByID)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
