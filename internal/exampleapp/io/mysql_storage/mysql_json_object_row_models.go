package mysql_storage

type mysqlJSONObjectRow struct {
	ID         int    `db:"id"`
	JSONObject string `db:"json_object"`
	SFObjectID string `db:"sf_object_id"`
	SchemaID   string `db:"schema_id"`
	UserID     string `db:"user_id"`
}
