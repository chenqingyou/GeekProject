package server

import (
	"fmt"
	"net/http"
)

type Server interface {
	// Route 设定一个路口，命中该路由会执行HandlerFunc的代码
	//pattern 路由的地址
	Route(pattern string, handlerFunc func(ctx *Context)) //创建server的时候就创建context
	// Start 启动服务器
	Start(address string) error
}

type SdkHttpServer struct {
	Name string
}

func (s *SdkHttpServer) Route(pattern string, handlerFunc func(ctx *Context)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) { //就这样写就行了不需要理解
		cxt := NewContext(writer, request) //希望暴露细节
		handlerFunc(cxt)
	}) //handlerFunc 传入一个方法实现它
}

func (s *SdkHttpServer) Start(address string) error {

	err := http.ListenAndServe(address, nil)
	if err != nil {
		return err
	}
	return nil
}

type commResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type signResponse struct {
	Email  string      `json:"email"`
	Name   string      `json:"name"`
	Passwd interface{} `json:"passwd"`
}

func SignUp(ctx *Context) {
	req := &signResponse{}
	err := ctx.ReadJson(req)
	if err != nil {
		ctx.BadRequestJson(err)
		return
	}
	reqs := &commResponse{Data: 123}
	err = ctx.WriteJson(http.StatusOK, reqs)
	if err != nil {
		fmt.Printf("写入响应失败：[%v]\n", err)
		return
	}
}

func NewSdkHttpServer(name string) *SdkHttpServer {
	return &SdkHttpServer{Name: name}
}
