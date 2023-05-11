/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type VerifyUser struct {
	// action that must be handled in module, must be \"verify_user\"
	Action string `json:"action"`
	// user email to verify with
	Email string `json:"email"`
	// user's id from identity
	UserId string `json:"user_id"`
}
