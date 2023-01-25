package gee

import (
	"fmt"
	"log"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || ((part[0] == '*' || part[0] == ':') && child.isWild) {
			return child
		} // prohibit: route overlap
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	wildnodes := make([]*node, 0)
	// put static route first.
	for _, child := range n.children {
		if child.part == part {
			nodes = append(nodes, child)
		} else if child.isWild {
			wildnodes = append(wildnodes, child)
		}
	}
	nodes = append(nodes, wildnodes...)
	return nodes
}

// below, actually second param is parts[height]
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { // only set pattern field in leave node
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil { // add a sibling layer
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		if child.isWild && len(n.children) > 0 {
			log.Panicf("%v: 存在路由冲突%v", part, n.children)
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// BFS
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { //reach leaf or wildcard character
		if n.pattern == "" { // as above, valid leaf with pattern field.
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
