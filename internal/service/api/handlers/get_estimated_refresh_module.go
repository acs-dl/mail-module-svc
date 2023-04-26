package handlers

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/acs/mail-module/internal/service/api/models"
	"gitlab.com/distributed_lab/ape"
)

func GetEstimatedRefreshModule(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, models.NewEstimatedTimeResponse(time.Duration(0).String()))
}
