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
				seg:    seg,
				param:  param,
				parent: currNode,
			}
			currNode.children = append(currNode.children, n)
			currNode = n
		}
	}
	if currNode.handler != nil {
		panic("two handlers in same path")
	}
	currNode.handler = h
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
	return currNode.handler
}
