package pubsub_publisher

type createEventJSON struct {
	ID           int    `json:"id"`
	SFObjectID   string `json:"sf_object_id"`
	EndUserID    string `json:"end_user_id"`
	CreationTime string `json:"creation_time"`
}

type deleteEventJSON struct {
	ID           int    `json:"id"`
	DeletionTime string `json:"deletion_time"`
}
