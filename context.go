package gun

import (
	"fmt"
	"net/http"
)

type Context struct {
	W        http.ResponseWriter
	R        *http.Request
	Method   string
	Path     string
	index    int
	handlers []HandlerFunc
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:      w,
		R:      r,
		Method: r.Method,
		Path:   r.URL.Path,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

// form参数查询
func (c *Context) PostForm(key string) string {
	return c.R.FormValue(key)
}

// query参数查询
func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

// 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

// string响应
func (c Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "文本/plain")
	if len(values) == 0 {
		c.W.Write([]byte(fmt.Sprintf(format)))
	} else {
		c.W.Write([]byte(fmt.Sprintf(format, values)))
	}
}
