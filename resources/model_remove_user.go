/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

// action to remove user from domain in mail
type RemoveUser struct {
	// action that must be handled in module, must be \"remove_user\"
	Action string `json:"action"`
	// user email to remove domain
	Email string `json:"email"`
}
