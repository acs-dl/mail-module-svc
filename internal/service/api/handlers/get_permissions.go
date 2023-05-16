package handlers

import (
	"net/http"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/internal/service/api/models"
	"github.com/acs-dl/mail-module-svc/internal/service/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetPermissionsRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	permissions := make([]data.Permission, 0)
	var totalCount int64 = 0

	permissionsQ := PermissionsQ(r).WithUsers()
	countPermissionsQ := PermissionsQ(r).CountWithUsers()

	userIds := make([]int64, 0)
	if request.UserId != nil {
		userIds = append(userIds, *request.UserId)
	}

	permissionsQ = permissionsQ.FilterByUserIds(userIds...)
	countPermissionsQ = countPermissionsQ.FilterByUserIds(userIds...)

	permissions, err = permissionsQ.Page(request.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get permissions")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	totalCount, err = countPermissionsQ.GetTotalCount()
	if err != nil {
		Log(r).WithError(err).Error("failed to get permissions total count")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.NewUserPermissionListResponse(permissions)
	response.Meta.TotalCount = totalCount
	response.Links = data.GetOffsetLinksForPGParams(r, request.OffsetPageParams)

	ape.Render(w, response)
}
