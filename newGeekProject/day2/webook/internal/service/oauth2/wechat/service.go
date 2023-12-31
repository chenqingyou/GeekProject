package wechat

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var redirectUrl = url.PathEscape("https://localhost.com//ouths2/wechat")

type ServiceWechatInterface interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}
type ServiceWechat struct {
	appId     string       `wire:"appid"`
	appSecret string       `wire:"appsecret"`
	client    *http.Client ``
}

func NewServiceWechat(appId string, appSecret string) ServiceWechatInterface {
	return &ServiceWechat{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient, //没有选择依赖注入，将来会变，但是现在没有变化
	}
}

func (sw *ServiceWechat) AuthUrl(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, sw.appId, redirectUrl, state), nil
}

func (sw *ServiceWechat) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
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

	//实现2将整个body都读取出来
	//Body, err := io.ReadAll(resp.Body)
	//err = json.Unmarshal(Body, &res)
	//if err != nil {
	//	return domain.WechatInfo{}, err
	//}

	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误信息%v[%v]\n", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenID:  res.OpenID,
		UnionID: res.UnionID,
	}, nil
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
