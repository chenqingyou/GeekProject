package domain

import (
	"time"
)

// UserDomain 领域对象，对应DDD中的聚合对象
type UserDomain struct {
	Id              int64     `json:"id,omitempty"`
	Email           string    `json:"email,omitempty"`
	Password        string    `json:"password,omitempty"`
	Nickname        string    `json:"nickname,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	Birthday        string    `json:"birthday,omitempty"`
	PersonalProfile string    `json:"personalProfile,omitempty"`
	Ctime           time.Time `json:"ctime,omitempty"`
}

type Result struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type WechatInfo struct {
}
