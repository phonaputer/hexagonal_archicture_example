package mysql_storage

import (
	"context"
	"database/sql"
	"errors"
	"examplemodule/internal/exampleapp/logic"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (r *UserStorage) Create(ctx context.Context, newUser *logic.NewUser) (int, error) {

	// Map from business logic data model to MySQL data model

	row := &mysqlUserRow{
		FirstName:    newUser.FirstName,
		LastName:     newUser.LastName,
		EmailAddress: newUser.EmailAddress,
	}

	// I/O. Execute MySQL query.

	const query = `INSERT INTO users (first_name, last_name, email_address) 
				   VALUES (:first_name, :last_name, :email_address)`

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

func (r *UserStorage) ExistsByEmailAddress(ctx context.Context, emailAddress string) (bool, error) {
	// No need to map from business logic data model to MySQL data model, because MySQL data model already
	// follows the BL model (emailAddress is a string in both places).

	const query = `SELECT EXISTS(SELECT * FROM users WHERE email_address = ?);`

	// I/O. Execute MySQL query.

	var result int
	err := r.db.QueryRowxContext(ctx, query, emailAddress).Scan(&result)
	if err != nil {
		return false, fmt.Errorf("select: %w", err)
	}

	// Map MySQL result to business logic data model & return result

	return result > 0, nil
}

func (r *UserStorage) GetByID(ctx context.Context, id int) (*logic.User, error) {
	// No need to map from business logic data model to MySQL data model, because MySQL data model already
	// follows BL data model (id is an integer in both places).

	// I/O. Execute MySQL query.

	const query = `SELECT * FROM users WHERE id = ?`

	var row mysqlUserRow

	err := r.db.QueryRowxContext(ctx, query, id).StructScan(&row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, logic.ErrUserNotFound // Map the MySQL error result to the business logic data model error
	}
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	// Map MySQL result to business logic data model & return result

	return &logic.User{
		ID:           row.ID,
		EmailAddress: row.EmailAddress,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
	}, nil
}

func (r *UserStorage) Delete(ctx context.Context, id int) error {
	// No need to map from business logic data model to MySQL data model, because MySQL data model
	// already follows BL data model (id is an integer in both places).

	// I/O. Execute MySQL query.

	const query = `DELETE FROM users WHERE id = ?`

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
		return logic.ErrUserNotFound
	}

	return nil
}
