package server

import (
	"encoding/json"
	"io"
	"net/http"
)

//为了简化SignUp里面的解析json和返回json的操作

type Context struct {
	W http.ResponseWriter //接口不需要用指针
	R *http.Request
}

func (c *Context) ReadJson(req interface{}) error {
	//处理json并且返回结构化之后的字符串
	r := c.R
	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(readBody, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) WriteJson(httpCode int, req interface{}) error {
	//处理json并且返回结构化之后的字符串
	c.W.WriteHeader(httpCode)

	marshal, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = c.W.Write(marshal)
	return err
}

func (c *Context) BadRequestJson(req interface{}) error {
	return c.WriteJson(http.StatusBadRequest, req)
}

func NewContext(writer http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: writer,
		R: r,
	}
}
