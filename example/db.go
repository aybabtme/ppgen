package example

//go:generate go run ../main.go nop -type UserDB -src db.go

type User struct {
	ID   string
	Name string
}

type UserDB interface {
	Create(name string) (*User, error)
	Get(id string) (u *User, ok bool, err error)
	Delete(*User) error
}
