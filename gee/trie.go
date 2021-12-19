package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	if n == nil {
		return nil
	}
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找（如果没有精确的匹配，也可以匹配child.isWild=1的那些）
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	nod := n
	for height < len(parts) {
		var child *node
		part := parts[height]
		child = nod.matchChild(part)
		if child == nil {
			//建立child节点
			child = &node{part: part, children: make([]*node, 0), isWild: part[0] == ':' || part[0] == '*'}
			nod.children = append(nod.children, child)
		}
		nod = child
		height++
	}
	nod.pattern = pattern
}

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
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

//遍历路由树，返回所有节点
func (n *node) travel(res *([]*node)) {
	if n == nil {
		return
	}
	if n.pattern != "" {
		*res = append(*res, n)
	}
	for _, child := range n.children {
		child.travel(res)
	}
}
