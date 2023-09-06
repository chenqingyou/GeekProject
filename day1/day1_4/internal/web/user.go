package web

import (
	"GeekProject/day1/day1_4/internal/domain"
	"GeekProject/day1/day1_4/internal/service"
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
	UserIdKey            = "userId"
	bizLogin             = "login"
)

type UserWebHandler struct {
	srv         *service.UserService
	passWordExp *regexp.Regexp
	emailExp    *regexp.Regexp
}

func NewUserWebHandler(svc *service.UserService) *UserWebHandler {
	return &UserWebHandler{
		srv:         svc,
		passWordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
	}
}

func (uwb *UserWebHandler) RegisterRoutes(engine *gin.Engine) {
	ug := engine.Group("/users")
	ug.POST("/signup", uwb.SignUp)
	ug.POST("/login", uwb.Login)
	ug.POST("/edit", uwb.Edit)
	ug.POST("/profile", uwb.Profile)
}

func (uwb *UserWebHandler) SignUp(ctx *gin.Context) {
	type userReqStruct struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req userReqStruct
	//校验，账号，密码，邮箱
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验邮箱
	emailExpOk, err := uwb.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error\n ")
		return
	}
	if !emailExpOk {
		ctx.String(http.StatusOK, "The registered email address format is incorrect\n ")
		return
	}
	//校验密码
	passWordExpOk, err := uwb.passWordExp.MatchString(req.Password)
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
	//服务层-数据层-数据库层
	err = uwb.srv.CreateUser(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmailRepSvc {
		ctx.String(http.StatusOK, "Mailbox conflict\n ")
		return
	}

	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error\n ")
		return
	}
	ctx.String(http.StatusOK, "Registered successfully\n ")
}

func (uwb *UserWebHandler) Login(ctx *gin.Context) {
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//进行密码比对
	userDetail, err := uwb.srv.LoginUser(ctx, domain.UserDomain{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserNotFoundSvc {
		ctx.String(http.StatusOK, "Mailbox conflict\n ")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "System error\n ")
		return
	}
	//设置session为数据库的id
	session := sessions.Default(ctx)
	session.Set(UserIdKey, userDetail.ID)
	session.Save()
	ctx.String(http.StatusOK, "Login successfully\n ")
}

func (uwb *UserWebHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname        string `json:"nickname"`
		Birthday        string `json:"birthday"`
		PersonalProfile string `json:"personalProfile"`
	}
	session := sessions.Default(ctx)
	userid := session.Get(UserIdKey).(int64)
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验字符串长度
	if utf8.RuneCountInString(req.Nickname) > 10 {
		ctx.JSON(http.StatusOK, domain.ResponseData{
			Info: "The nickname length cannot exceed 10",
			Code: 200001,
			Data: "",
		})
		return
	}
	//校验字符串长度
	if utf8.RuneCountInString(req.PersonalProfile) > 300 {
		ctx.JSON(http.StatusOK, domain.ResponseData{
			Info: "The personalProfile length cannot exceed 300",
			Code: 200001,
			Data: "",
		})
		return
	}
	// 使用 time 包中的 Parse 函数解析日期字符串
	_, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, domain.ResponseData{
			Info: "Birthday is formatted incorrectly, for example:2006-01-01",
			Code: 200001,
			Data: "",
		})
		return
	}
	//调用service层的编辑接口
	err = uwb.srv.EditUser(ctx, domain.UserDomain{
		ID:              userid,
		Nickname:        req.Nickname,
		Birthday:        req.Birthday,
		PersonalProfile: req.PersonalProfile,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, domain.ResponseData{
			Info: "Description Failed to edit user information",
			Code: 200001,
			Data: "",
		})
		return
	}
	ctx.JSON(http.StatusOK, domain.ResponseData{
		Info: "User information is successfully edited",
		Code: 200001,
		Data: "",
	})
}

func (uwb *UserWebHandler) Profile(ctx *gin.Context) {

}
