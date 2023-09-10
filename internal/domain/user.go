package domain

type User struct {
	id       uint64
	userName string
}

func NewUser(id uint64, userName string) *User {
	return &User{
		id:       id,
		userName: userName,
	}
}

func (u *User) ID() uint64 {
	return u.id
}

func (u *User) SetID(id uint64) {
	u.id = id
}

func (u *User) UserName() string {
	return u.userName
}

func (u *User) SetUserName(userName string) {
	u.userName = userName
}
