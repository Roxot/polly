package httpserver

import "polly/database"

type UserWithToken struct {
	Id          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Token       string `json:"token"`
	DisplayName string `json:"display_name"`
	DeviceType  int    `json:"device_type"`
}

func UserToUserWithToken(user database.User) UserWithToken {
	uwt := UserWithToken{}
	uwt.Id = user.Id
	uwt.PhoneNumber = user.PhoneNumber
	uwt.Token = user.Token
	uwt.DisplayName = user.DisplayName
	uwt.DeviceType = user.DeviceType
	return uwt
}
