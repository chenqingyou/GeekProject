package ioc

import "GeekProject/newGeekProject/day2/webook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.ServiceWechatInterface {
	return wechat.NewServiceWechat("appId")
}
