package gee

import (
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Write      http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int

	Params map[string]string

	// 中间件
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Write:  w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Write.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Write.Header().Set(key, value)
}

func (c *Context) String(code int, str string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Write.Write([]byte(str))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Write)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Write, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Write.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.Write.Write([]byte(html))
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}
