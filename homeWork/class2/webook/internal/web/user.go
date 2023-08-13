package web

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/service"
	"encoding/json"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

	fmt.Printf("%v\n", req)
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
	ctx.String(http.StatusOK, "Registered successfully\n ")
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
	sess.Set(userIdKey, uLoginMeg.Id)
	//设置之后需要刷新
	sess.Save()
	ctx.String(http.StatusOK, "Login successful\n")
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
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验字符串长度
	if utf8.RuneCountInString(req.Nickname) > 10 {
		ctx.String(http.StatusOK, "The nickname length cannot exceed 10\n ")
		return
	}
	//校验字符串长度
	if utf8.RuneCountInString(req.PersonalProfile) > 300 {
		ctx.String(http.StatusOK, "The personalProfile length cannot exceed 300\n ")
		return
	}
	// 使用 time 包中的 Parse 函数解析日期字符串
	_, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "Birthday is formatted incorrectly, for example:2006-01-02\n ")
		return
	}
	//调用服务层的登录接口
	err = u.svc.Edit(ctx, domain.UserDomain{
		Id:              id,
		Nickname:        req.Nickname,
		Birthday:        req.Birthday,
		PersonalProfile: req.PersonalProfile,
	})
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "Mailbox does not exist \n")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error")
		return
	}
	ctx.String(http.StatusOK, "edit successfully\n ")
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	type EditReq struct {
		Email string `json:"email"`
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//调用服务层的登录接口
	userDetail, err := u.svc.Profile(ctx, id)
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
