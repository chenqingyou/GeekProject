package domain

import "time"

type UserDomain struct {
	ID              int64
	Email           string
	Password        string
	Nickname        string `json:"nickname"`
	Birthday        string `json:"birthday"`
	PersonalProfile string `json:"personalProfile"`
	Ctime           time.Time
}

type ResponseData struct {
	Info string `json:"info"`
	Code int    `json:"code"`
	Data string `json:"data"`
}
