package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
)

func GetUserRolesMap(w http.ResponseWriter, r *http.Request) {
	result := newModuleRolesResponse()

	result.Data.Attributes["super_admin"] = "write"
	result.Data.Attributes["admin"] = "write"
	result.Data.Attributes["user"] = "read"

	ape.Render(w, result)
}
