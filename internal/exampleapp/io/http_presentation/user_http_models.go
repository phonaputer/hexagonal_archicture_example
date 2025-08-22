package http_presentation

type userCreateRequest struct {
	FirstName    string
	EmailAddress string
	LastName     string
}

type userResponse struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstName"`
	EmailAddress string `json:"emailAddress"`
	LastName     string `json:"lastName"`
}

type userDeleteRequest struct {
	ID int
}

type userGetByIDRequest struct {
	ID int
}
