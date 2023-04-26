package handlers

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/acs/mail-module/internal/service/api/models"
	"gitlab.com/distributed_lab/acs/mail-module/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetEstimatedRefreshSubmodule(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewRefreshSubmoduleRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, models.NewEstimatedTimeResponse(time.Duration(0).String()))
}
