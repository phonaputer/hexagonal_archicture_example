package mysql_storage

type mysqlUserRow struct {
	ID           int    `db:"id"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	EmailAddress string `db:"email_address"`
}
