package models

import (
	"strconv"

	"github.com/acs-dl/mail-module-svc/internal/data"
	"github.com/acs-dl/mail-module-svc/resources"
)

func NewUserInfoModel(user data.User, id int) resources.UserInfo {
	return resources.UserInfo{
		Key: resources.Key{
			ID:   strconv.Itoa(id),
			Type: resources.USER,
		},
		Attributes: resources.UserInfoAttributes{
			Email: user.Email,
			Name:  user.Name,
		},
	}
}

func NewUserInfoList(users []data.User, offset uint64) []resources.UserInfo {
	result := make([]resources.UserInfo, len(users))
	for i, user := range users {
		result[i] = NewUserInfoModel(user, i+int(offset))
	}
	return result
}

func NewUserInfoListResponse(users []data.User, offset uint64) resources.UserInfoListResponse {
	return resources.UserInfoListResponse{
		Data: NewUserInfoList(users, offset),
	}
}
