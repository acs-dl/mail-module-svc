package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Users interface {
	New() Users

	Upsert(user User) error
	Delete() error
	Select() ([]User, error)
	Get() (*User, error)

	FilterById(id *int64) Users
	FilterByMailIds(emailIds ...string) Users
	FilterByLowerTime(time time.Time) Users
	SearchBy(search string) Users

	Count() Users
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams) Users
}

type User struct {
	Id        *int64    `json:"-" db:"id" structs:"id,omitempty"`
	Email     string    `json:"email" db:"email" structs:"email"`
	Name      string    `json:"name" db:"name" structs:"name"`
	MailId    string    `json:"mail_id" db:"mail_id" structs:"mail_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" structs:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" structs:"-"`
	Submodule *string   `json:"-" db:"-" structs:"-"`
}

type UnverifiedUser struct {
	CreatedAt time.Time `json:"created_at"`
	Module    string    `json:"module"`
	Submodule string    `json:"submodule"`
	ModuleId  string    `json:"module_id"`
	Email     *string   `json:"email,omitempty"`
	Name      *string   `json:"name,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	Username  *string   `json:"username,omitempty"`
}
