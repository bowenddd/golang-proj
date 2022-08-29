package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node //roots用于保存每个方法的trie树的根节点，便于路由
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if strings.HasPrefix(item, "*") {
			break
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	n, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
		n = r.roots[method]
	}
	n.insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	parts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(parts, 0)
	if n == nil {
		return nil, nil
	}
	nParts := parsePattern(n.pattern)
	for i, part := range nParts {
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = parts[i]
		}
		if strings.HasPrefix(part, "*") {
			params[part[1:]] = strings.Join(parts[i:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Req.Method, c.Req.URL.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, r.handlers[key])
		// r.handlers[key](c)
	} else {
		c.handlers = append(c.handlers, HandlerFunc(func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND! %s\n", c.Path)
		}))
	}
	// 按顺序执行handler
	c.Next()
}
