package web

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/service"
	"encoding/json"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	userIdKey            = "userId"
	bizLogin             = "login"
)

// 确保UserHandler实现了handler接口
var _ Handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserServiceInterface
	passWordExp *regexp.Regexp
	emailExp    *regexp.Regexp
	codeSvc     service.CodeServiceInterface
	jwtHandler  //用组合的方式，不使用指针
}

func NewUserHandler(svc service.UserServiceInterface, codeSvc service.CodeServiceInterface) *UserHandler {
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		passWordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		jwtHandler:  NewJwtHandler(),
	}
}

// RegisterRoutesV1 使用传入分组的方式使用路由
func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	//注册
	ug.POST("/signup", u.SignUp)
	//登录
	ug.POST("/login", u.Longin)
	//修改
	ug.POST("/edit", u.Edit)
	//查询
	ug.POST("/profile", u.Profile)
}

func (u *UserHandler) RegisterRoutesCt(server *gin.Engine) {
	//可以使用分组注册路由的方法
	ug := server.Group("/users")
	//注册
	ug.POST("/signup", u.SignUp)
	//登录
	ug.POST("/login", u.Longin)
	ug.POST("/loginJwt", u.LonginJwt)
	//修改
	ug.POST("/edit", u.Edit)
	//查询
	ug.POST("/profile", u.Profile)
	//发送验证码
	ug.POST("/loginSms/code", u.SendLoginSMSCode)
	//校验验证码
	ug.POST("/loginSms", u.LoginSMS)
	ug.POST("/refreshToken", u.RefreshToken)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	//bind 方法会根据Content-Type 来解析你的数据到req里面
	//解析错了就会返回一个404的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验邮箱
	emailExpOk, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error\n ")
		return
	}
	if !emailExpOk {
		ctx.String(http.StatusOK, "The registered email address format is incorrect\n ")
		return
	}
	//校验密码
	passWordExpOk, err := u.passWordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error\n ")
		return
	}
	if !passWordExpOk {
		ctx.String(http.StatusOK, "The password must be larger than 8 characters and contain special characters\n ")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "The two passwords are different\n  ")
		return
	}

	//下面是服务端的注册用户
	err = u.svc.SignUp(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "Mailbox conflict")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	ctx.String(http.StatusOK, "Registered successfully")
}

//RefreshToken 可以同时刷新长短token，用redis来记录是否有效，即使refresh_token是一次性的
//参考登录校验部分，比较User-Agent

func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	//只有这个接口，拿出来的是refresh_token，其他的都是短token
	refreshToken := ExtractToken(ctx)
	var rc RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return u.rtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.SetJWTToken(ctx, rc.Uid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 0,
		Msg:  "刷新成功",
		Data: nil,
	})
}

func (u *UserHandler) Longin(ctx *gin.Context) {
	type SignUpReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//调用服务层的登录接口
	uLoginMeg, err := u.svc.Login(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "The account or password is incorrect\n \n")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	//这里登录成功了
	sess := sessions.Default(ctx)
	//放在sess里的值
	sess.Set("userId", uLoginMeg.Id)
	sess.Options(sessions.Options{
		//Path:     "",
		//Domain:   "",
		MaxAge:   1 * 60,
		Secure:   false, //要求只能使用https
		HttpOnly: false, //当设置了 HttpOnly 标志，这意味着该 cookie 不能通过客户端
		// JavaScript 进行访问。这有助于保护 cookie 不被恶意脚本获取和操纵，尤其是在跨站脚本 (XSS) 攻击的环境下。
		//SameSite: 0,
	})
	//设置之后需要刷新
	sess.Save()
	ctx.String(http.StatusOK, "Login successful\n")
	return
}

func (u *UserHandler) LonginJwt(ctx *gin.Context) {
	type SignUpReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//调用服务层的登录接口
	uLoginMeg, err := u.svc.Login(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "The account or password is incorrect\n \n")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	err = u.SetJWTToken(ctx, uLoginMeg.Id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	err = u.setRefreshToken(ctx, uLoginMeg.Id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	ctx.String(http.StatusOK, "Login successful\n")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	//放在sess里的值
	sess.Options(sessions.Options{
		//Path:     "",
		//Domain:   "",
		MaxAge:   -1,    //退出登录
		Secure:   false, //要求只能使用https
		HttpOnly: false, //当设置了 HttpOnly 标志，这意味着该 cookie 不能通过客户端
		// JavaScript 进行访问。这有助于保护 cookie 不被恶意脚本获取和操纵，尤其是在跨站脚本 (XSS) 攻击的环境下。
		//SameSite: 0,
	})
	//设置之后需要刷新
	sess.Save()
	ctx.String(http.StatusOK, "exitLogin successful\n")
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	//先判读类型
	claims, ok := c.(*TokenClaims)
	if !ok {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	fmt.Printf("claims[%v]\n", claims)
	ctx.String(http.StatusOK, "这是你的 Profile ")
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	type EditReq struct {
		Email string `json:"email"`
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*TokenClaims)
	if !ok {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	fmt.Printf("claims[%v]\n", claims)
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//调用服务层的登录接口
	userDetail, err := u.svc.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	//需要将返回值反解码
	userDetailMarshal, err := json.Marshal(userDetail)
	if err != nil {
		return
	}
	ctx.String(http.StatusOK, string(userDetailMarshal))
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	const biz = "login"
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "输入有误",
			Data: nil,
		})
	}
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, domain.Result{
			Msg: "验证码发送成功",
		})
	case service.ErrSetCodeFrequently:
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "验证码发送太频繁",
			Data: nil,
		})
	default:
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "验证码发送失败",
			Data: nil,
		})
	}

}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	const biz = "login"
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "手机号码输入有误",
			Data: nil,
		})
	}
	if req.Code == "" {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "验证码输入有误",
			Data: nil,
		})
	}
	judgment, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		if err == service.ErrVerifyCodeFrequently {
			ctx.JSON(http.StatusOK, domain.Result{
				Code: 5,
				Msg:  "验证次数过多",
			})
			return
		}
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !judgment {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 这边要怎么办呢？
	// 从哪来？
	if err = u.SetJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = u.setRefreshToken(ctx, user.Id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}

	ctx.JSON(http.StatusOK, domain.Result{
		Msg: "验证码校验通过",
	})

}
