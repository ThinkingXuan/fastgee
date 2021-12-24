package fastgee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context struct
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middleware
	handlers []HandlerFunc
	index    int
}

// newContext create a context
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next() 这个设计非常巧妙：共有两个地方需要调用next()：
// 中间件和最终处理的handle
// 如果有一个logger的全局中间件 + 一个普通的get请求， 他们都会回调用next方法。
// handers会有两个hanldefunc(loggerFunc, getFunc)
// loggerFunc如果存在next方法，回使用loggerFunc中的next去调用getFunc,
// 能调用原因： context只有一份，index是共享变量，loggerFunc的next调用后，回将index增加到和len(c.handlers)的长度。

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm Post Form request
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query query request
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status update context statusCode and header
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Fail(code int, message string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": message})
}
