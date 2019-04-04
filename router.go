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
		params   map[string][]string    // map{method: [param, param]}
	}
)

// NewRouter create an new Router instance
func NewRouter() *Router {
	r := &Router{
		root: &node{
			seg:      "",
			handlers: make(map[string]HandlerFunc),
			params:   make(map[string][]string),
		},
	}
	return r
}

// Register register a route rule, use trie tree
func (r *Router) Register(method string, path string, h HandlerFunc) {
	// split path by '/'
	segs := strings.Split(path, "/")
	// put every item into tree
	currNode := r.root
	seg, param := "", ""
	var params []string
	for i, v := range segs {
		if i == 0 {
			continue
		}
		if len(v) > 1 && v[:1] == ":" {
			seg = "*"
			param = v[1:]
			params = append(params, param)
		} else {
			seg = v
			param = ""
		}
		isMatched := false
		for _, cn := range currNode.children {
			if cn.seg == seg {
				isMatched = true
				currNode = cn
				break
			}
		}
		if !isMatched {
			n := &node{
				seg:      seg,
				param:    param,
				parent:   currNode,
				handlers: make(map[string]HandlerFunc),
				params:   make(map[string][]string),
			}
			currNode.children = append(currNode.children, n)
			currNode = n
		}
	}
	currNode.setHandler(method, h)
	currNode.setParams(method, params)
}

// Detect detect `HandlerFunc` correspond with given `path`, param saved in Context
func (r *Router) Detect(path string, ctx *Context) HandlerFunc {
	segs := strings.Split(path, "/")
	currNode := r.root
	var pvalues []string
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
			pvalues = append(pvalues, v)
		} else if !isMatched {
			return nil
		}
	}
	// save param
	method := ctx.request.Method
	for k, v := range MakeMapBySlices(currNode.getParams(method), pvalues) {
		ctx.SetParam(k, v)
	}
	return currNode.getHandler(method)
}

func (n *node) fmtMethod(method string) string {
	if method == "" {
		method = "*"
	}
	method = strings.ToUpper(method)
	return method
}

// setHandler set handler with method, * means accept any method
func (n *node) setHandler(method string, handler HandlerFunc) {
	method = n.fmtMethod(method)
	_, ok := n.handlers[method]
	if ok {
		panic(fmt.Sprintf("%s handler conflict", method))
	}
	n.handlers[method] = handler
}

// getHandler get handler by method
func (n *node) getHandler(method string) HandlerFunc {
	method = n.fmtMethod(method)
	h, ok := n.handlers[method]
	if !ok {
		h = n.handlers["*"]
	}
	return h
}

// setParams set params with method, * means accept any method
func (n *node) setParams(method string, params []string) {
	method = n.fmtMethod(method)
	_, ok := n.params[method]
	if ok {
		panic(fmt.Sprintf("%s handler conflict", method))
	}
	n.params[method] = params
}

// getParams get params by method
func (n *node) getParams(method string) []string {
	method = n.fmtMethod(method)
	h, ok := n.params[method]
	if !ok {
		h = n.params["*"]
	}
	return h
}

func MakeMapBySlices(keys []string, vals []string) map[string]string {
	res := make(map[string]string)
	keys_len := len(keys)
	vals_len := len(vals)
	if keys_len != vals_len {
		return nil
	}
	for i, key := range keys {
		res[key] = vals[i]
	}
	return res
}
