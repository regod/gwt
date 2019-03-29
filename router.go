package gwt

import (
	"strings"
)

type (
	Router struct {
		root *node
	}
	node struct {
		seg      string
		param    string
		handler  HandlerFunc
		parent   *node
		children []*node
	}
)

// NewRouter create an new Router instance
func NewRouter() *Router {
	r := &Router{
		root: &node{
			seg: "",
		},
	}
	return r
}

// Register register a route rule, use trie tree
func (r *Router) Register(path string, h HandlerFunc) {
	// path clean, remove '/' at the end if exist
	path = strings.TrimSuffix(path, "/")
	// split path by '/'
	segs := strings.Split(path, "/")
	// put every item into tree
	curr_node := r.root
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
		for _, cn := range curr_node.children {
			if cn.seg == seg {
				isMatched = true
				curr_node = cn
				if seg == "*" && cn.param != param {
					panic("different param")
				}
				break
			}
		}
		if !isMatched {
			n := &node{
				seg:    seg,
				param:  param,
				parent: curr_node,
			}
			curr_node.children = append(curr_node.children, n)
			curr_node = n
		}
	}
	if curr_node.handler != nil {
		panic("two handlers in same path")
	}
	curr_node.handler = h
}

// Detect detect `HandlerFunc` correspond with given `path`, param saved in Context
func (r *Router) Detect(path string, ctx *Context) HandlerFunc {
	path = strings.TrimSuffix(path, "/")
	segs := strings.Split(path, "/")
	curr_node := r.root
	for i, v := range segs {
		if i == 0 {
			continue
		}

		isMatched := false
		var wild_node *node
		for _, cn := range curr_node.children {
			if cn.seg == "*" {
				wild_node = cn
				continue
			}
			if cn.seg == v {
				isMatched = true
				curr_node = cn
				break
			}
		}
		if !isMatched && wild_node != nil {
			curr_node = wild_node
			// save param
			ctx.param[curr_node.param] = v
		} else if !isMatched {
			return nil
		}
	}
	return curr_node.handler
}
