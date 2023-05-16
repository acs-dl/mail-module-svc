/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserPermissionAttributes struct {
	// domain name
	Link string `json:"link"`
	// user id from module
	ModuleId string `json:"module_id"`
	// domain name
	Path string `json:"path"`
	// user id from identity
	UserId *int64 `json:"user_id,omitempty"`
	// email from domain
	Username string `json:"username"`
}
