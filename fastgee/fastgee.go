package fastgee

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of server http
type Engine struct {
	//router map[string]HandlerFunc
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

/*
RouterGroup 持有Engine指针，Engine又继承RouterGroup。
原理： Engine 需要用有管理路由的功能， RouterGroup也要有管理路由的功能，
      防止功能重叠，所以让他们相互持有。Engine拥有RouterGroup所有功能。
	  RouterGroup持有持有Engine指针也可以进行路由控制。

*/

// New is the constructor of fastgee Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//// addRoute
//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//	engine.router.addRoute(method, pattern, handler)
//}
//
//// GET defines the method to add GET request
//func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//	engine.addRoute("GET", pattern, handler)
//}
//
//// POST defines the method to add POST request
//func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//	engine.addRoute("POST", pattern, handler)
//}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP implement interface Handler's method ServeHTTP
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

type RouterGroup struct {
	prefix string
	// support middleware
	middleware []HandlerFunc
	parent     *RouterGroup
	// all groups share a engin instance
	engine *Engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
