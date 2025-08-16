package model

type Message struct {
	ID       string
	UserID   string
	Payload  string
	Language string
}

type User struct {
	ID        string
	Name      string
	AccountID string
}

type Account struct {
	ID      string
	Balance int64
}

type Profile struct {
	User
	Account
}
