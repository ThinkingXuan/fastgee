package fastgee

import (
	"net/http"
	"strings"
)

// router 路由控制
// 主要的作用: newRouter()  addRoute() getRoute() handle()
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

// parsePattern 解析路由到字符数组里面
func parsePattern(pattern string) []string {

	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 如果第一个根路由的第一个字符为 *,直接停止赋值。
				break
			}
		}
	}
	return parts
}

// addRoute
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 以method为key创建roots
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	// 构造前缀树
	r.roots[method].insert(pattern, parts, 0)

	// 以 method + "-" + pattern 为key绑定路由函数
	r.handlers[key] = handler
}

// getRoute
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 解析 请求的path
	searchParts := parsePattern(path)

	params := make(map[string]string)

	// 判断次请求方法是否存在，不存在直接结束
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 搜索待匹配的路由
	n := root.search(searchParts, 0)

	// 解析匹配到节点的pattern，主要是为了处理:和* 两种情况，解析出来他们的参数。
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				// 例如 匹配路由index/:name 请求为index/:hello  params存储为，key=name value=hello
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil

}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)

	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
