package data

import (
	"time"
)

const (
	ModuleName = "mail"
)

type ModuleRequest struct {
	ID            string    `db:"id" structs:"id"`
	UserID        int64     `db:"user_id" structs:"user_id"`
	Module        string    `db:"module" structs:"module"`
	Payload       string    `db:"payload" structs:"payload"`
	CreatedAt     time.Time `db:"created_at" structs:"created_at"`
	RequestStatus string    `db:"request_status" structs:"request_status"`
	Error         string    `db:"error" structs:"error"`
}

type ModulePayload struct {
	RequestId string `json:"request_id"`
	UserId    string `json:"user_id"`
	Action    string `json:"action"`

	//other fields that are required for module
	Link      string   `json:"link"`
	Links     []string `json:"links"`
	Email     string   `json:"username"`
	Phone     string   `json:"phone"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
}

type UnverifiedPayload struct {
	Action string           `json:"action"`
	Users  []UnverifiedUser `json:"users"`
}
