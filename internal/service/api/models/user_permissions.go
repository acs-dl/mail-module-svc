package models

import (
	"strconv"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/resources"
)

func NewUserPermissionModel(permission data.Permission, counter int) resources.UserPermission {
	result := resources.UserPermission{
		Key: resources.Key{
			ID:   strconv.Itoa(counter),
			Type: resources.USER_PERMISSION,
		},
		Attributes: resources.UserPermissionAttributes{
			ModuleId: permission.User.MailId,
			UserId:   permission.Id,
			Link:     permission.Link,
			Path:     permission.Link,
			Username: permission.Email,
		},
	}

	return result
}

func NewUserPermissionList(permissions []data.Permission) []resources.UserPermission {
	result := make([]resources.UserPermission, len(permissions))
	for i, permission := range permissions {
		result[i] = NewUserPermissionModel(permission, i)
	}
	return result
}

func NewUserPermissionListResponse(permissions []data.Permission) UserPermissionListResponse {
	return UserPermissionListResponse{
		Data: NewUserPermissionList(permissions),
	}
}

type UserPermissionListResponse struct {
	Meta  Meta                       `json:"meta"`
	Data  []resources.UserPermission `json:"data"`
	Links *resources.Links           `json:"links"`
}
