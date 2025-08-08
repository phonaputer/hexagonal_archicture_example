package http_presentation

type jsonObjectCreateRequest struct {
	SchemaID   string
	SFObjectID string
	EndUserID  string
	JSONObject string
}

type jsonObjectResponse struct {
	ID         int    `json:"id"`
	SchemaID   string `json:"schemaId"`
	SFObjectID string `json:"sfObjectId"`
	EndUserID  string `json:"endUserId"`
	JSONObject string `json:"jsonObject"`
}

type jsonObjectDeleteRequest struct {
	ID int
}

type jsonObjectGetByIDRequest struct {
	ID int
}
