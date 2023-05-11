package handlers

import (
	"net/http"

	"github.com/acs-dl/mail-module-svc/internal/service/api/models"
	"github.com/acs-dl/mail-module-svc/internal/service/api/requests"
	"github.com/acs-dl/mail-module-svc/internal/service/googleApi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUsersRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	email := ""
	if request.Email != nil {
		email = *request.Email
	}

	users, err := UsersQ(r).SearchBy(email).Select()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to select users from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if len(users) != 0 {
		ape.Render(w, models.NewUserInfoListResponse(users, 0))
		return
	}

	users, err = googleApi.GoogleClientInstance(ParentContext(r.Context())).SearchByUsersFromApi(email)
	if err != nil {
		Log(r).WithError(err).Infof("failed to get users from api by `%s`", email)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewUserInfoListResponse(users, 0))
}
