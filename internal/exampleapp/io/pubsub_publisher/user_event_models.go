package pubsub_publisher

type createUserEventJSON struct {
	ID           int    `json:"id"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	CreationTime string `json:"creation_time"`
}

type deleteUserEventJSON struct {
	ID           int    `json:"id"`
	DeletionTime string `json:"deletion_time"`
}
