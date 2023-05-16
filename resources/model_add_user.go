/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type AddUser struct {
	// action that must be handled in module, must be \"add_user\"
	Action string `json:"action"`
	// user email to add domain
	Email string `json:"email"`
	// first name of the user
	FirstName string `json:"first_name"`
	// last name of the user
	LastName string `json:"last_name"`
	// link where module has to add user
	Link string `json:"link"`
	// user's id from identity
	UserId string `json:"user_id"`
}
