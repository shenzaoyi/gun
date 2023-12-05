package gun

import "fmt"

// import (
//
//	"strings"
//
// )
//
//	type Router struct {
//		roots    map[string]*node       //	一个方法一棵树
//		handlers map[string]HandlerFunc //	path : HandlerFunc
//	}
//
//	func NewRouter() *Router {
//		return &Router{
//			roots:    make(map[string]*node),
//			handlers: make(map[string]HandlerFunc),
//		}
//	}
//
// //	解析路径
//
//	func ParsePattern(pattern string) []string {
//		vs := strings.Split(pattern, "/")
//		parts := make([]string, 0)
//		for _, v := range vs {
//			if v != "" {
//				if v != "*" {
//					parts = append(parts, v)
//				}
//				break
//			}
//		}
//		return parts
//	}
//
//	func (r *Router) AddRouter(method, pattern string, handler HandlerFunc) {
//		if _, ok := r.roots[method]; !ok {
//			r.roots[method] = &node{}
//		}
//		parts := ParsePattern(pattern)
//		r.roots[method].Insert(pattern, parts, 0)
//		path := method + "-" + pattern
//		r.handlers[path] = handler
//	}
//
//	func (r *Router) GetRouter(method, pattern string) {
//		if _, ok := r.roots[method]; !ok {
//
//		}
//		parts := ParsePattern(pattern)
//		r.roots[method].Search(parts, 0)
//	}
var r Router = Router{"GET",
	&node{}}
var routerTree []*Router = []*Router{
	&r,
}

type Router struct {
	Method string
	Root   *node
}

func (r *Router) handle(c *Context) {
	//	根据请求路径获取对应函数
	fmt.Println("来到了r.handle")
	value := r.Root.GetValue(c.Path)
	if value.handlers != nil {
		//for _, v := range value.handlers {
		//	v(c)
		//}
		c.handlers = append(c.handlers, value.handlers...)
		fmt.Println("长度的为")
		fmt.Println(len(c.handlers))
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(200, "404 NOT FOUND")
		})
	}
	c.Next()
}
