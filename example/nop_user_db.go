// Code generated by ppgen (github.com/aybabtme/ppgen).
//
// command:
// 		ppgen nop -type UserDB -src db.go
//
// DO NOT EDIT!

package example

// NopUserDB returns a UserDB that does nothing.
func NopUserDB() UserDB { return nopUserDB{} }

type nopUserDB struct{}

func (nopUserDB) Create(_ string) (out0 *User, out1 error)   { return out0, out1 }
func (nopUserDB) Get(_ string) (u *User, ok bool, err error) { return u, ok, err }
func (nopUserDB) Delete(_ *User) (out0 error)                { return out0 }
