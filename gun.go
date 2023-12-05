package gun

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(c *Context)
type HandlersChain []HandlerFunc

type Engine struct {
	*RouterGroup
	Routers []Router
	Groups  []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	MiddleWares []HandlerFunc
	engine      *Engine
	parent      *RouterGroup
}

func (r *RouterGroup) Group(prefix string, middleware ...HandlerFunc) *RouterGroup {
	rg := RouterGroup{
		prefix:      r.prefix + prefix,
		MiddleWares: middleware,
		engine:      r.engine,
		parent:      r,
	}
	r.engine.Groups = append(r.engine.Groups, &rg)
	return &rg
}

// RG首先获取完整的中间件，拼接上handlers，注册到路由树上
func (r *RouterGroup) addRoute(method, path string, handlers HandlersChain) {
	methodRoot := getMethodTree(method, r.engine)
	//	获取中间件,拼接
	handlers = append(handlers, r.MiddleWares...)
	methodRoot.Root.AddRouter(path, handlers)
}

func (r *RouterGroup) GET(pattern string, handlers ...HandlerFunc) {
	path := r.prefix + pattern
	r.addRoute("GET", path, handlers)
}

func New() *Engine {
	engine := Engine{Routers: make([]Router, 0, 9)}
	engine.RouterGroup = &RouterGroup{engine: &engine}
	engine.Groups = make([]*RouterGroup, 0)
	return &engine
}

func (e *Engine) Run(addr string) {
	fmt.Println(fmt.Sprintf("服务启动起来了，跑在了%s 端口，祝开发顺利", addr))
	err := http.ListenAndServe(addr, e)
	if err != nil {
		fmt.Println(err)
	}
}

func getMethodTree(method string, e *Engine) *Router {
	methodRoot := &Router{
		Method: "",
		Root:   new(node),
	}
	for _, v := range e.Routers {
		if v.Method == "GET" {
			methodRoot = &v
		}
	}
	if methodRoot.Method == "" {
		methodRoot.Method = "GET"
		e.Routers = append(e.Routers, *methodRoot)
	}
	return methodRoot
}

func (e *Engine) GET(pattern string, handler ...HandlerFunc) {
	//	首先找到当前方法树
	methodRoot := getMethodTree("GET", e)
	//	将当前处理链加到trie
	methodRoot.Root.AddRouter(pattern, handler)
}

func handle(e *Engine, c *Context) error {
	//	首先根据方法获取方法树
	fmt.Println("来到了handle")
	router := Router{
		Method: "",
		Root:   new(node),
	}
	method := c.Method
	for _, v := range e.Routers {
		if v.Method == method {
			router = v
		}
	}
	if router.Method == "" {
		return fmt.Errorf("获取方法树失败")
	}
	router.handle(c)
	return nil
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//	构建Context
	c := newContext(w, r)
	fmt.Println("来到了ServeHTTP")
	//c := Context{
	//	W:        nil,
	//	R:        nil,
	//	Method:   "GET",
	//	Path:     "/search",
	//	index:    -1,
	//	handlers: nil,
	//}
	//	router
	err := handle(e, c)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (rg *RouterGroup) Use(handlers ...HandlerFunc) {
	//
	rg.MiddleWares = append(rg.MiddleWares, handlers...)

}
