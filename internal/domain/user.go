package domain

type User struct {
	id       string
	userName string
}

func NewUser(id, userName string) *User {
	return &User{
		id:       id,
		userName: userName,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) SetID(id string) {
	u.id = id
}

func (u *User) UserName() string {
	return u.userName
}

func (u *User) SetUserName(userName string) {
	u.userName = userName
}
