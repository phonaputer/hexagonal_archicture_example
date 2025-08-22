package logic

type NewUser struct {
	EmailAddress string
	FirstName    string
	LastName     string
}

type User struct {
	ID           int
	EmailAddress string
	FirstName    string
	LastName     string
}
