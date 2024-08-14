package models

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	VoipId      string `json:"voipId"`
	Active      bool   `json:"active"`
	PushToken   string `json:"pushToken"`
}
