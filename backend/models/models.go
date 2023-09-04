package models

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Chat struct {
	IsGroupChat bool   `json:"isGroupChat"`
	Users       []User `json:"users"`
	ID          string `json:"_id"`
	ChatName    string `json:"chatName"`
	GroupAdmin  *User  `json:"groupAdmin,omitempty"`
}
