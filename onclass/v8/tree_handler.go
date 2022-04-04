package main

import (
	"net/http"
	"strings"
)

type HandlerBasedTree struct{
	root *Node
}

type Node struct {
	path string
	children []*Node
	handler handlerFunc
}

func (h *HandlerBasedTree) ServerHTTP(c *Context) {
	handler, found := h.findRouter(c.R.URL.Path)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)
}

func (h *HandlerBasedTree) findRouter(path string) (handlerFunc, bool) {
	// 去除头尾可能有的/，然后按照/切割成段
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := h.root
	for _, p := range paths {
		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, found := h.findMathChild(cur, p)
		if !found {
			return nil, false
		}
		cur = matchChild
	}
	// 到这里，应该是找完了
	if cur.handler == nil {
		// 到达这里是因为这种场景
		// 比如说你注册了 /user/profile
		// 然后你访问 /user
		return nil, false
	}
	return cur.handler, true
}

func (h *HandlerBasedTree) Route(method string, pattern string, handleFun handlerFunc) {
	patt := strings.Trim(pattern,"/")	//处理路径的前后斜杆， /user/friends  /user/friends/  user/friends/
	paths := strings.Split(patt, "/")

	cur := h.root
	for index, path := range paths{
		mathchild, ok := h.findMathChild(cur, path)
		if ok{
			cur = mathchild
		} else {
			h.createSubTree(cur, paths[index:], handleFun)
		}
	}
}

func (h *HandlerBasedTree) findMathChild(root *Node, path string) (*Node, bool){
	for _, child := range root.children{
		if child.path == path{
			return child, true
		}
	}
	return nil, false
}

//
func (h *HandlerBasedTree) createSubTree(root *Node, paths []string, handlerFun handlerFunc) {
	cur := root
	for _, path := range paths{
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handlerFun
}

func newNode(path string) *Node {
	return &Node{
		path: path,
		children: make([]*Node, 0, 4),
	}
}