package radix

import (
	"bytes"
	"sort"
)

// 节点类型定义
type Node struct {
	// 节点存储的数据
	Value interface{}

	// 边缘节点
	edges []*edge

	// 节点优先级
	priority int

	// 节点深度
	depth int
}

// 获取节点深度
func (n *Node) Depth() int {
	return n.depth
}

// 判断该节点是否为叶节点
func (n *Node) IsLeaf() bool {
	length := len(n.edges)
	// 兼容二进制树
	if length == 2 {
		return n.edges[0] == nil && n.edges[1] == nil
	}
	return length == 0
}

// 获取节点优先级
func (n *Node) Priority() int {
	return n.priority
}

// 增加二进制树节点
func (n *Node) addBinary(label string, v interface{}) (nn int) {
	for i := range label {
		for j := uint8(8); j > 0; j-- {
			bbit := bit(j, label[i])
			done := i == len(label)-1 && j == 1
			if e := n.edges[bbit]; e != nil {
				if done {
					e.n.Value = v
					return
				}
				goto next
			}
			n.edges[bbit] = &edge{
				n: &Node{
					depth: n.depth + 1,
					edges: make([]*edge, 2),
				},
			}
			if done {
				n.edges[bbit].n.Value = v
			}
			nn++
		next:
			n = n.edges[bbit].n
		}
	}
	return nn
}

func (n *Node) clone() *Node {
	c := *n
	c.incrDepth()
	return &c
}

func (n *Node) delBinary(label string) int {
	var (
		ref *edge
		del int
	)
	for i := range label {
		for j := uint8(8); j > 0; j-- {
			bbit := bit(j, label[i])
			done := i == len(label)-1 && j == 1
			if e := n.edges[bbit]; e != nil {
				del++
				if done && e.n.IsLeaf() { // 只在节点为叶节点时删除，否则会破坏树
					ref.n.edges = make([]*edge, 2) // 从最后一个有值的节点重置边
					return del
				}
				ref = e
				n = e.n
				continue
			}
			return 0
		}
	}
	return 0
}

func (n *Node) getBinary(label string) *Node {
	for i := range label {
		for j := uint8(8); j > 0; j-- {
			bbit := bit(j, label[i])
			done := i == len(label)-1 && j == 1
			if e := n.edges[bbit]; e != nil {
				if done {
					return e.n
				}
				n = e.n
				continue
			}
			return nil
		}
	}
	return nil
}

func (n *Node) incrDepth() {
	n.depth++
	for _, e := range n.edges {
		e.n.incrDepth()
	}
}

// 递归排序节点及其子节点
func (n *Node) sort(st SortingTechnique) {
	s := &sorter{
		n:  n,
		st: st,
	}
	sort.Sort(s)
	for _, e := range n.edges {
		e.n.sort(st)
	}
}

func (n *Node) writeTo(bd *builder) {
	for i, e := range n.edges {
		e.writeTo(bd, []bool{i == len(n.edges)-1})
	}
}

func (n *Node) writeToBinary(bd *builder, buf, aux *bytes.Buffer) {
	prefix := aux.Bytes()
	length := len(prefix)
	aux1, aux2 := make([]byte, length), make([]byte, length)
	copy(aux1, prefix)
	copy(aux2, prefix)
	auxs := []*bytes.Buffer{
		bytes.NewBuffer(aux1),
		bytes.NewBuffer(aux2),
	}
	for i, e := range n.edges {
		if e != nil {
			bit := byte('0')
			if i == 1 {
				bit = '1'
			}
			auxs[i].WriteByte(bit)
			if e.n != nil {
				if e.n.Value != nil {
					bd.Write(prefix)
					bd.WriteByte(bit) // holds only one value
					isLeaf := e.n.IsLeaf()
					if isLeaf {
						bd.WriteString(bd.colors[colorGreen].Wrap(" 🍂"))
					}
					bd.WriteString(bd.colors[colorMagenta].Wrapf(" → %#v\n", e.n.Value))
				}
				e.n.writeToBinary(bd, buf, auxs[i])
			}
		}
	}
}

func bit(i uint8, c byte) uint8 {
	if 1<<(i-1)&c > 0 {
		return 1
	}
	return 0
}
