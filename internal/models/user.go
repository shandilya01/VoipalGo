package models

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Incantation string `json:"incantation"`
	Active      bool   `json:"active"`
	PushToken   string `json:"pushToken"`
}
