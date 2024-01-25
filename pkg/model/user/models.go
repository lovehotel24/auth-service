package user

import (
	"time"
	"unsafe"
)

// dbUser represent the structure we need for moving data
// between the app and the database.
type dbUser struct {
	ID           string    `db:"user_id"`
	Name         string    `db:"name"`
	Phone        string    `db:"phone"`
	Role         string    `db:"role"`
	PasswordHash []byte    `db:"password_hash"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

// User represents an individual user.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	PasswordHash []byte    `json:"-"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string `json:"name" validate:"required"`
	Phone           string `json:"phone" validate:"required,email"`
	Role            string `json:"role" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}

// UpdateUser defines what information to provided to modify an existing
// User.
type UpdateUser struct {
	Name            *string `json:"name"`
	Phone           *string `json:"phone" validate:"omitempty,email"`
	Role            *string `json:"role"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}

func toDBUser(u User) dbUser {
	pdb := (*dbUser)(unsafe.Pointer(&u))
	return *pdb
}

func toUser(db dbUser) User {
	pu := (*User)(unsafe.Pointer(&db))
	return *pu
}

func toUserSlice(dbUsers []dbUser) []User {
	users := make([]User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		users[i] = toUser(dbUsr)
	}
	return users
}
