package domain

import mysms "GeekProject/homeWork/class2/webook/internal/service/sms"

type SmsDomain struct {
	Tpl     string
	NameArg []mysms.NameArg
	Numbers []string
}
