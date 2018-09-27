package models

type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Address   string `json:"address"`
	Gender    string `json:"gender"`
	Timestamp int64  `json:"timestamp"`
}

type Response struct {
	Values []string `json:"values"`
}
