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
	svc         *service.UserService
	passWordExp *regexp.Regexp
	emailExp    *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		passWordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
	}
}

// RegisterRoutesV1 使用传入分组的方式使用路由
//func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
//	//注册
//	ug.POST("/signup", u.SignUp)
//	//登录
//	ug.POST("/login", u.Longin)
//	//修改
//	ug.POST("/edit", u.Edit)
//	//查询
//	ug.POST("/profile", u.Profile)
//}

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
		return
	}
	//调用服务层的登录接口
	uLoginMeg, err := u.svc.Login(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 0,
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

	//使用jwt
	//privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) //密钥
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, domain.Result{
	//		Code: 0,
	//		Msg:  " fmt.Sprintf(\"err [%v]\", err)",
	//		Data: nil,
	//	})
	//	return
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
		ID:               uLoginMeg.Id,
		UserAgent:        ctx.Request.UserAgent(),
	})
	signedString, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.JSON(http.StatusOK, domain.Result{
			Code: 5,
			Msg:  "SignedString err",
			Data: nil,
		})
		return
	}
	log.Printf("x-jwt-tokn[%v]\n", signedString)
	ctx.Header("x-jwt-tokn", signedString)
	ctx.Header("x-jwt-tokn", signedString)
	ctx.JSON(http.StatusOK, domain.Result{
		Code: 0,
		Msg:  "Login successful",
		Data: nil,
	})
	return
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
		Email string `json:"email"`
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
		Msg:  "查询成功",
		Data: userDetail,
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	ID        int64
	UserAgent string
}
