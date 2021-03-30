package fastgee

import "strings"

type node struct {
	pattern  string  // 待匹配的路由,例如：/index/:user/*file
	part     string  // 路由中的一部分   例如：index
	children []*node // 子节点      所有的下级路由节点
	isWild   bool    // 不是精确匹配  如/:user /*file  为true
}

// matchChild 第一匹配成功的节点，用户插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 匹配到 part，或者不是精确匹配 /: /* 结束函数
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	// 匹配到 part，或者不是精确匹配 /: /*，添加进入nodes放回
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 插入路由，
func (n *node) insert(pattern string, parts []string, height int) {

	// 匹配到最后一层，赋值pattern，其他层数不赋值
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 是否存在匹配part的孩子节点
	part := parts[height]
	child := n.matchChild(part)

	// 不存在，则创建并加入到children，注意 : 和 *的情况，isWild赋值为true
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 递归插入路由下一层
	child.insert(pattern, parts, height+1)
}
// search
func (n *node) search(parts []string, height int) *node {

	// 搜索到最后一层或者出现*，并且pattern不是空时，返回该节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 查找此路由节点下所有可能的孩子
	part := parts[height]
	children := n.matchChildren(part)

	// 遍历所有孩子，并递归查找。
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
