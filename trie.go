package gun

import "fmt"

//type node struct {
//	pattern  string
//	part string
//	children []*node
//	isWild   bool
//}
//
//func (n *node)matchChild(part string)*node{
//	for _, child := range n.children{
//		if child.part == part || child.isWild{
//			return child
//		}
//	}
//	return nil
//}
//
//func (n *node)matchChildren(part string)[]*node{
//	nodes := make([]*node,0)
//	for _, child := range n.children{
//		if child.part == part || child.isWild{
//			nodes = append(nodes,child)
//		}
//	}
//	return nodes
//}
//
////	api/v1/user
//func (n *node) Insert(pattern string, parts []string, height int) {
//	if len(parts) == height{
//		n.pattern = pattern
//		//	出去的条件，parts的高度等于height
//		return
//	}
//	part := parts[height]
//	child := n.matchChild(part)
//	if child == nil{
//		child = &node{
//			pattern:  "",
//			part:     part,
//			children: nil,
//			isWild:   parts[0] == ":" || parts[0] == "*",
//		}
//		n.children = append(n.children,child)
//	}
//	child.Insert(pattern,parts,height+1)
//}
//
//
//func (n *node)Search(parts []string, height int)*node{
//	if len(parts) == height{
//		if n.pattern == ""{
//			return nil
//		}
//		return n
//	}
//	part := parts[height]
//	children := n.matchChildren(part)
//	for _, child := range children{
//		result := child.Search(parts,height+1)
//		if result != nil{
//			return result
//		}
//	}
//	return nil
//}

type node struct {
	relativePath string
	absolutePath string
	indices      string //	索引
	children     []*node
	priority     int
	handlers     HandlersChain
}

type nodeValue struct {
	absolutePath string
	handlers     HandlersChain
}

func (nv nodeValue) Get(ctx *Context) {
	for _, v := range nv.handlers {
		v(ctx)
	}
}

// 添加

func (n *node) InsertChild(fullPath string, handlers HandlersChain) {
	n.relativePath = fullPath
	n.absolutePath = fullPath
	n.handlers = handlers
}

// search/see : se + arch
//	求最长前缀

func (n *node) LongestPrefix(path string) int {
	i := 0
	lenth := 0
	if len(path) > len(n.relativePath) {
		lenth = len(n.relativePath)
	} else {
		lenth = len(path)
	}
	for j := 0; j < lenth; j++ {
		if path[j] == n.relativePath[j] {
			i += 1
		} else {
			return i
		}
	}
	return i
}

// 根据前缀得到属于第几个孩子
func (n *node) GetPrefixChild(a string) int {
	for k, v := range n.children {
		if v.relativePath[0:1] == a {
			return k
		}
	}
	return -1
}

// 添加

func (n *node) AddRouter(path string, handlers HandlersChain) {
	//	如果当前为根结点，则直接加到根节点
	if len(n.relativePath) == 0 && len(n.children) == 0 {
		n.InsertChild(path, handlers)
		fmt.Println(n.relativePath)
		return
	}
	//	如果不是
	//	先求最大公共前缀
	i := n.LongestPrefix(path)
	child := node{
		relativePath: n.relativePath[i:],
		absolutePath: "",
		indices:      n.indices,
		children:     n.children,
		priority:     n.priority,
		handlers:     n.handlers,
	}
	n.children = []*node{&child}
	n.relativePath = n.relativePath[:i]
	child.absolutePath = n.relativePath + child.relativePath
	n.absolutePath = ""
	n.indices = child.relativePath[:1] //	设置孩子首字母为indices第一个
	n.handlers = nil
	path = path[i:]
	c := path[:1]
	for i := 0; i < len(n.indices); i++ { //	检索是否存在仍可以利用的公共前缀
		if n.indices == c {
			i := n.GetPrefixChild(c)
			n = n.children[i]
			n.AddRouter(path, handlers)
		}
	}
	//	如果无公共前缀，直接添加
	abPath := n.relativePath + path
	childNew := node{
		relativePath: path,
		absolutePath: abPath,
		indices:      "",
		children:     nil,
		priority:     0,
		handlers:     handlers,
	}
	n.children = append(n.children, &childNew)
	n.indices += childNew.relativePath[0:1]
}

//		根据根节点展示路径, 只是测试用
//
//func Show(n *node) {
//	if n.absolutePath != "" {
//		fmt.Println(n.absolutePath)
//	} else {
//		if n.indices == "" {
//			return
//		} else {
//			for _, v := range n.children {
//				Show(v)
//			}
//		}
//
//	}
//}

func (n *node) GetValue(path string) (value nodeValue) {
	prefix := n.relativePath
	if len(path) > len(prefix) {
		if path[:len(prefix)] == prefix {
			path = path[len(prefix):]
			idxc := path[0]
			for _, c := range []byte(n.indices) {
				if c == idxc {
					i := n.GetPrefixChild(string([]byte{c}))
					n = n.children[i]
					value = n.GetValue(path)
				}
			}
		}
	}
	if len(path) == len(prefix) {
		if value.handlers = n.handlers; value.handlers != nil {
			value.absolutePath = n.absolutePath
			return value
		}
	}
	return
}
