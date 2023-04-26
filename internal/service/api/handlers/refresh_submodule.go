package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/acs/mail-module/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RefreshSubmodule(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewRefreshSubmoduleRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, http.StatusAccepted)
}
