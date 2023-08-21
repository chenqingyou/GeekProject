package web

import (
	"GeekProject/newGeekProject/day2/webook/internal/domain"
	"GeekProject/newGeekProject/day2/webook/internal/service"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
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
	ug.POST("/loginJwt", u.LonginJwt)
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

	//使用 ecdsa.GenerateKey 生成了一个 ECDSA 密钥对，并将私钥用于 JWT 的签名
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"info": fmt.Sprintf("err [%v]", err)})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES512, TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), //过期时间
			NotBefore: nil,
		},
		Uid:       uLoginMeg.Id,
		UserAgent: ctx.Request.UserAgent(),
	})
	signedString, err := token.SignedString(privateKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"info": fmt.Sprintf("err [%v]", err)})
		return
	}
	fmt.Printf("signedString [%v]\n", signedString)
	fmt.Printf("uLoginMeg [%v]\n", uLoginMeg)
	ctx.Header("x-jwt-token", signedString)
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
	ctx.String(http.StatusOK, "这是你的 Profile ")
}

// TokenClaims 实现jwt的接口
type TokenClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
