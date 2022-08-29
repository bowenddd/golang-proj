package gee

import "strings"

type node struct {
	pattern  string  // 总的路径 e.g. /users/bwdeng/info
	part     string  // 当前节点的地址 e.g. info
	children []*node //子节点
	isWild   bool    // 是否精确匹配，用于动态路由。true表示非精确匹配（动态路由）
}

// 查找第一个匹配的子节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, ch := range n.children {
		if ch.part == part || ch.isWild {
			return ch
		}
	}
	return nil
}

// 查找所有匹配的子节点，用于查找
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, ch := range n.children {
		if ch.part == part || ch.isWild {
			children = append(children, ch)
		}
	}
	return children
}

// 插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*")}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 查找
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		node := child.search(parts, height+1)
		if node != nil {
			return node
		}
	}
	return nil
}
