package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Permissions interface {
	New() Permissions

	Upsert(permission Permission) error
	Delete() error
	Select() ([]Permission, error)
	Get() (*Permission, error)

	FilterByMailIds(mailIds ...string) Permissions
	FilterByGreaterTime(time time.Time) Permissions
	FilterByLowerTime(time time.Time) Permissions
	FilterByLinks(links ...string) Permissions

	WithUsers() Permissions
	FilterByUserIds(userIds ...int64) Permissions

	Count() Permissions
	CountWithUsers() Permissions
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams) Permissions
}

type Permission struct {
	RequestId string    `json:"request_id" db:"request_id" structs:"request_id"`
	MailId    string    `json:"mail_id" db:"mail_id" structs:"mail_id"`
	Link      string    `json:"link" db:"link" structs:"link"`
	CreatedAt time.Time `json:"created_at" db:"created_at" structs:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" structs:"-"`
	*User     `structs:",omitempty"`
}
