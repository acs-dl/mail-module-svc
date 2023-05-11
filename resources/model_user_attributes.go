/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"time"
)

type UserAttributes struct {
	// timestamp without timezone when user was created
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// module name
	Module string `json:"module"`
	// submodule name
	Submodule *string `json:"submodule,omitempty"`
	// user id from identity module, if user is not verified - null
	UserId *int64 `json:"user_id,omitempty"`
	// email from mail
	Username *string `json:"username,omitempty"`
}
