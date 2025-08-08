package mysql_storage

import (
	"context"
	"database/sql"
	"errors"
	"examplemodule/internal/exampleapp/logic"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type JSONObjectStorage struct {
	db *sqlx.DB
}

func NewJSONObjectStorage(db *sqlx.DB) *JSONObjectStorage {
	return &JSONObjectStorage{db: db}
}

func (r *JSONObjectStorage) Create(ctx context.Context, jsonObject *logic.NewJSONObject) (int, error) {

	// Map from business logic data model to MySQL data model

	row := &mysqlJSONObjectRow{
		JSONObject: jsonObject.Object,
		SFObjectID: jsonObject.SFObjectID,
		SchemaID:   jsonObject.SchemaID,
		UserID:     jsonObject.UserID,
	}

	// I/O. Execute MySQL query.

	const query = `INSERT INTO json_objects (json_object, sf_object_id, schema_id, user_id) 
				   VALUES (:json_object, :sf_object_id, :schema_id, :user_id)`

	res, err := r.db.NamedExecContext(ctx, query, row)
	if err != nil {
		return 0, fmt.Errorf("insert: %w", err)
	}

	// Map MySQL result to business logic data model & return result

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return int(id), err
}

func (r *JSONObjectStorage) ExistsBySFObjectID(ctx context.Context, sfObjectID string) (bool, error) {
	// No need to map from business logic data model to MySQL data model, because MySQL data model already
	// follows the BL model (sfObjectID is a string in both places).

	const query = `SELECT EXISTS(SELECT * FROM json_objects WHERE sf_object_id = ?);`

	// I/O. Execute MySQL query.

	var result int
	err := r.db.QueryRowxContext(ctx, query, sfObjectID).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("select: %w", err)
	}

	// Map MySQL result to business logic data model & return result

	return result > 0, nil
}

func (r *JSONObjectStorage) GetByID(ctx context.Context, id int) (*logic.JSONObject, error) {
	// No need to map from business logic data model to MySQL data model, because MySQL data model already
	// follows BL data model (id is an integer in both places).

	// I/O. Execute MySQL query.

	const query = `SELECT * FROM json_objects WHERE id = ?`

	var row mysqlJSONObjectRow

	err := r.db.QueryRowxContext(ctx, query, id).StructScan(&row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logic.ErrObjectNotFound // Map the MySQL error result to the business logic data model error
	}
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	// Map MySQL result to business logic data model & return result

	return &logic.JSONObject{
		ID:         row.ID,
		Object:     row.JSONObject,
		SFObjectID: row.SFObjectID,
		SchemaID:   row.SchemaID,
		UserID:     row.UserID,
	}, nil
}

func (r *JSONObjectStorage) Delete(ctx context.Context, id int) error {
	// No need to map from business logic data model to MySQL data model, because MySQL data model
	// already follows BL data model (id is an integer in both places).

	// I/O. Execute MySQL query.

	const query = `DELETE FROM json_objects WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("affected rows: %w", err)
	}

	// Map MySQL result to business logic data model.
	// (In this case, by mapping "0 affected rows" to the business logic error).

	if affectedRows < 1 {
		return logic.ErrObjectNotFound
	}

	return nil
}
