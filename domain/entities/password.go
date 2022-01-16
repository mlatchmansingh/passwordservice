package entities

type Password struct {
	Id        int64  `json:"id"`
	Password  string `json:"password"`
	Converted bool
}

func NewPassword(plain string) Password {
	return Password{
		Id:        NewID(),
		Password:  plain,
		Converted: false,
	}
}
