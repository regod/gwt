package gwt

import (
	"fmt"
	"strings"
)

type (
	Router struct {
		root *node
	}
	node struct {
		seg      string
		param    string
		parent   *node
		children []*node
		handlers map[string]HandlerFunc // map{method: handler}
	}
)

// NewRouter create an new Router instance
func NewRouter() *Router {
	r := &Router{
		root: &node{
			seg:      "",
			handlers: make(map[string]HandlerFunc),
		},
	}
	return r
}

// Register register a route rule, use trie tree
func (r *Router) Register(method string, path string, h HandlerFunc) {
	// path clean, remove '/' at the end if exist
	path = strings.TrimSuffix(path, "/")
	// split path by '/'
	segs := strings.Split(path, "/")
	// put every item into tree
	currNode := r.root
	seg, param := "", ""
	for i, v := range segs {
		if i == 0 {
			continue
		}
		if v[:1] == ":" {
			seg = "*"
			param = v[1:]
		} else {
			seg = v
			param = ""
		}
		isMatched := false
		for _, cn := range currNode.children {
			if cn.seg == seg {
				isMatched = true
				currNode = cn
				if seg == "*" && cn.param != param {
					panic("different param")
				}
				break
			}
		}
		if !isMatched {
			n := &node{
				seg:      seg,
				param:    param,
				parent:   currNode,
				handlers: make(map[string]HandlerFunc),
			}
			currNode.children = append(currNode.children, n)
			currNode = n
		}
	}
	currNode.setHandler(method, h)
}

// Detect detect `HandlerFunc` correspond with given `path`, param saved in Context
func (r *Router) Detect(path string, ctx *Context) HandlerFunc {
	path = strings.TrimSuffix(path, "/")
	segs := strings.Split(path, "/")
	currNode := r.root
	for i, v := range segs {
		if i == 0 {
			continue
		}

		isMatched := false
		var wildNode *node
		for _, cn := range currNode.children {
			if cn.seg == "*" {
				wildNode = cn
				continue
			}
			if cn.seg == v {
				isMatched = true
				currNode = cn
				break
			}
		}
		if !isMatched && wildNode != nil {
			currNode = wildNode
			// save param
			ctx.SetParam(currNode.param, v)
		} else if !isMatched {
			return nil
		}
	}
	return currNode.getHandler(ctx.request.Method)
}

// setHandler set handler with method, * means accept any method
func (n *node) setHandler(method string, handler HandlerFunc) {
	if method == "" {
		method = "*"
	}
	method = strings.ToUpper(method)
	_, ok := n.handlers[method]
	if ok {
		panic(fmt.Sprintf("%s handler conflict", method))
	}
	n.handlers[method] = handler
}

// getHandler get handler by method
func (n *node) getHandler(method string) HandlerFunc {
	method = strings.ToUpper(method)
	h, ok := n.handlers[method]
	if !ok {
		h = n.handlers["*"]
	}
	return h
}
