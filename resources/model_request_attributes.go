/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RequestAttributes struct {
	// Module to grant permission
	Module string `json:"module"`
	// Already built payload to grant permission <br><br> -> \"add_user\" = action to add user in module<br> -> \"verify_user\" = action to verify user in module (connect user id from identity with module info)<br> -> \"get_users\" = action to get users with their permissions from module<br> -> \"delete_user\" = action to delete user from module<br> -> \"remove_user\" = action to remove user from submodule<br>
	Payload json.RawMessage `json:"payload"`
}
