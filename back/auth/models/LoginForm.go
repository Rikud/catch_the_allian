package models

type LoginForm struct {
	Email    string
	Password string
}

func (f LoginForm) GetPassword() string {
	return f.Password
}

func (f LoginForm) GetEmail() string {
	return f.Email
}
