package web

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	userIdKey            = "userId"
	bizLogin             = "login"
)

type UserHandler struct {
	svc         service.UserServiceInterface
	svcCode     service.CodeServiceInterface
	passWordExp *regexp.Regexp
	emailExp    *regexp.Regexp
}

func NewUserHandler(svc service.UserServiceInterface, svcCode service.CodeServiceInterface) *UserHandler {
	return &UserHandler{
		svc:         svc,
		svcCode:     svcCode,
		passWordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutesCt(server *gin.Engine) {
	//可以使用分组注册路由的方法
	ug := server.Group("/users")
	//注册
	ug.POST("/signup", u.SignUp)
	//登录
	ug.POST("/login", u.Longin)
	//修改
	ug.POST("/edit", u.Edit)
	//查询
	ug.POST("/profile", u.Profile)
	ug.POST("/loginSms/code", u.SendLoginSMSCode)
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
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	if !emailExpOk {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "The registered email address format is incorrect",
			Data: nil,
		})
		return
	}
	//校验密码
	passWordExpOk, err := u.passWordExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 0,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	if !passWordExpOk {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "The password must be larger than 8 characters and contain special characters",
			Data: nil,
		})
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "The two passwords are different",
			Data: nil,
		})
		return
	}

	fmt.Printf("%v\n", req)
	//下面是服务端的注册用户
	err = u.svc.SignUp(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "Mailbox conflict",
			Data: nil,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 0,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 0,
		Msg:  "Registered successfully",
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
		fmt.Printf("err[%v]\n", err)
		return
	}
	//调用服务层的登录接口
	uLoginMeg, err := u.svc.Login(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 4,
			Msg:  "The account or password is incorrect",
			Data: nil,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	if u.JWTToken(ctx, uLoginMeg, err) {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 0,
		Msg:  "Login successful",
		Data: nil,
	})
	return
}

func (u *UserHandler) JWTToken(ctx *gin.Context, uLoginMeg domain.UserDomain, err error) bool {
	fmt.Printf("uLoginMeg[%v]\n", uLoginMeg)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
		ID:               uLoginMeg.Id,
		UserAgent:        ctx.Request.UserAgent(),
	})
	if token == nil {
		return true
	}
	signedString, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		return true
	}
	log.Printf("x-jwt-tokn[%v]\n", signedString)
	ctx.Header("x-jwt-tokn", signedString)
	return false
}

/*
允许用户补充基本个人信息，包括：
昵称：字符串，你需要考虑允许的长度。
生日：前端输入为 1992-01-01 这种字符串。
个人简介：一段文本，你需要考虑允许的长度。
尝试校验这些输入，并且返回准确的信息。
修改 /users/profile 接口，确保这些信息也能输出到前端。*/

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname        string `json:"nickname"`
		Birthday        string `json:"birthday"`
		PersonalProfile string `json:"personalProfile"`
	}
	//sess := sessions.Default(ctx)
	//id := sess.Get(userIdKey).(int64)
	//使用jwt
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	fmt.Printf("claims[%v]\n", claims)

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验字符串长度
	if utf8.RuneCountInString(req.Nickname) > 10 {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "The nickname length cannot exceed 10",
			Data: nil,
		})
		return
	}
	//校验字符串长度
	if utf8.RuneCountInString(req.PersonalProfile) > 300 {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "The personalProfile length cannot exceed 300",
			Data: nil,
		})
		return
	}
	// 使用 time 包中的 Parse 函数解析日期字符串
	_, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "Birthday is formatted incorrectly, for example:2006-01-02",
			Data: nil,
		})
		return
	}
	//调用服务层的登录接口
	err = u.svc.Edit(ctx, domain.UserDomain{
		Id:              claims.ID,
		Nickname:        req.Nickname,
		Birthday:        req.Birthday,
		PersonalProfile: req.PersonalProfile,
	})
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
			Msg:  "Mailbox does not exist",
			Data: nil,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	//成功将数据缓存到本地
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 0,
		Msg:  "edit successfully",
		Data: nil,
	})
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	type EditReq struct {
		Email           string `json:"email,omitempty"`
		Nickname        string `json:"nickname,omitempty"`
		Birthday        string `json:"birthday,omitempty"`
		PersonalProfile string `json:"personalProfile,omitempty"`
	}
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	fmt.Printf("claims[%v]\n", claims)
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//调用服务层的登录接口
	userDetail, err := u.svc.Profile(ctx, claims.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.Result{
			Code: 5,
			Msg:  "System error",
			Data: nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.Result{
		Msg: "查询成功",
		Data: EditReq{
			Email:           userDetail.Email,
			Nickname:        userDetail.Nickname,
			Birthday:        userDetail.Birthday,
			PersonalProfile: userDetail.PersonalProfile,
		}, //最好不要将领域对象直接暴露出去，需要重新定义一个结构体
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	ID        int64
	UserAgent string
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
	err := u.svcCode.Send(ctx, biz, req.Phone)
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
	judgment, err := u.svcCode.Verify(ctx, biz, req.Phone, req.Code)
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
	_, err = u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 这边要怎么办呢？
	// 从哪来？
	//if err = u.SetJWTToken(ctx, user.Id); err != nil {
	//	ctx.JSON(http.StatusOK, domain.Result{
	//		Code: 5,
	//		Msg:  "系统错误",
	//	})
	//	return
	//}

	ctx.JSON(http.StatusOK, domain.Result{
		Msg: "验证码校验通过",
	})

}
