package wechat

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"net/url"
)

var redirectUrl = url.PathEscape("https://localhost.com//ouths2/wechat")

type ServiceWechatInterface interface {
	AuthUrl(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code, state string) (domain.WechatInfo, error)
}
type serviceWechat struct {
	appId     string
	appSecret string
	client    *http.Client
}

func NewServiceWechat(appId, appSecret string) ServiceWechatInterface {
	return &serviceWechat{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient, //没有选择依赖注入，将来会变，但是现在没有变化
	}
}

func (sw *serviceWechat) AuthUrl(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, sw.appId, redirectUrl, state), nil
}

func (sw *serviceWechat) VerifyCode(ctx context.Context, code, state string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, sw.appId, sw.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := sw.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误信息%v[%v]\n", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{}, nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenID  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}
