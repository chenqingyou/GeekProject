package domain

import "time"

// UserDomain 领域对象，对应DDD中的聚合对象
type UserDomain struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}
