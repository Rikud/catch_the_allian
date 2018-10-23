package models

type User struct {
	Id 		 int 	`json:"-"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score 	 int `json:"score, string"`
}

func (user *User) GetIdPoint() *int {
	return &user.Id
}

func (user *User) GetId() int {
	return user.Id
}

func (user *User) GetEmailPoint() *string {
	return &(user.Email)
}

func (user *User) GetPasswordPoint() *string {
	return &(user.Password)
}

func (user *User) GetUsernamePoint() *string {
	return &(user.Username)
}

func (user *User) GetAvatarPoint() *string {
	return &(user.Avatar)
}

func (user *User) GetEmail() string {
	return user.Email
}

func (user *User) GetPassword() string {
	return user.Password
}

func (user *User) GetUsername() string {
	return user.Username
}

func (user *User) GetAvatar() string {
	return user.Avatar
}

func (user *User) SetId(id int) {
	user.Id = id
}

func (user *User) SetEmail(email string) {
	user.Email = email
}

func (user *User) SetPassword(password string) {
	user.Password = password
}

func (user *User) SetUsername(username string) {
	user.Username = username
}

func (user *User) SetAvatar(avatar string) {
	user.Avatar = avatar
}

func (user *User) SetScore(score int) {
	user.Score = score
}