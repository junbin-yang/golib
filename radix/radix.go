package radix

import (
	"bytes"
	"github.com/gbrlsnchs/color"
	"strings"
	"sync"
)

const (
	// Tsafe 是否启线程安全
	Tsafe = 1 << iota
	// Tdebug 添加更多描述信息.
	Tdebug
	// Tbinary 默认使用二进制树而不是前缀树
	Tbinary

	Tnocolor
)

// 树结构定义
type Tree struct {
	// 根节点
	root *Node

	// 节点总数
	length int

	// 节点字节长度
	size int

	// 是否需要锁
	safe bool

	// 二进制树或前缀树
	binary bool

	// 占位符，用于匹配动态参数
	placeholder byte

	// 分隔符，用于分割字符串，如url里的/
	delim byte

	mu *sync.RWMutex
	bd *builder
}

// 初始化一个二叉树
func BinaryRoot() *Tree {
	return New(Tbinary)
}

// 初始化一个动态前缀树
func PrefixRoot() *Tree {
	t := New(0)
	t.SetBoundaries(':', '/')
	return t
}

// 初始化一个根节点
func New(flags int) *Tree {
	tr := &Tree{
		root:   &Node{},
		length: 1,
	}
	if flags&Tbinary > 0 {
		tr.binary = true
		tr.root.edges = make([]*edge, 2)
	}
	if flags&Tsafe > 0 {
		tr.mu = &sync.RWMutex{}
		tr.safe = true
	}
	tr.bd = &builder{
		Builder: &strings.Builder{},
		debug:   flags&Tdebug > 0,
	}

	tr.bd.colors[colorRed] = color.New(color.CodeFgRed)
	tr.bd.colors[colorGreen] = color.New(color.CodeFgGreen)
	tr.bd.colors[colorMagenta] = color.New(color.CodeFgMagenta)
	tr.bd.colors[colorBold] = color.New(color.CodeBold)
	for _, c := range tr.bd.colors {
		c.SetDisabled(flags&Tnocolor > 0)
	}

	return tr
}

// 添加一个新节点
func (tr *Tree) Add(label string, v interface{}) {
	// 不允许有空接口和空的节点数据.
	if label == "" || v == nil {
		return
	}
	if tr.safe {
		defer tr.mu.Unlock()
		tr.mu.Lock()
	}
	tnode := tr.root
	if tr.binary {
		nn := tnode.addBinary(label, v)
		tr.length += nn
		tr.size += nn / 8
		return
	}
	for {
		var next *edge
		var slice string
		for _, edge := range tnode.edges {
			var found int
			slice = edge.label
			for i := range slice {
				if i < len(label) && slice[i] == label[i] {
					found++
					continue
				}
				break
			}
			if found > 0 {
				label = label[found:]
				slice = slice[found:]
				next = edge
				break
			}
		}
		if next != nil {
			tnode = next.n
			tnode.priority++
			// 匹配完整的单词
			if len(label) == 0 {
				// 如果节点标签与边缘节点的标签一致，则替换。
				// Example:
				// 	(root) -> tnode("tomato", v1)
				// 	becomes
				// 	(root) -> tnode("tomato", v2)
				if len(slice) == 0 {
					tnode.Value = v
					return
				}
				// 标签是边缘节点标签的前缀
				// Example:
				// 	(root) -> tnode("tomato", v1)
				// 	then add "tom"
				// 	(root) -> ("tom", v2) -> ("ato", v1)
				next.label = next.label[:len(next.label)-len(slice)]
				c := tnode.clone()
				c.priority--
				tnode.edges = []*edge{
					&edge{
						label: slice,
						n:     c,
					},
				}
				tnode.Value = v
				tr.length++
				return
			}
			// 添加一个新节点，但将其父节点分解为前缀和将剩余的切片作为新的边缘节点。
			// Example:
			// 	(root) -> ("tomato", v1)
			// 	then add "tornado"
			// 	(root) -> ("to", nil) -> ("mato", v1)
			// 	                      +> ("rnado", v2)
			if len(slice) > 0 {
				c := tnode.clone()
				c.priority--
				tnode.edges = []*edge{
					&edge{ // 复制后缀到新节点
						label: slice,
						n:     c,
					},
					&edge{ // 新节点
						label: label,
						n: &Node{
							Value:    v,
							depth:    tnode.depth + 1,
							priority: 1,
						},
					},
				}
				next.label = next.label[:len(next.label)-len(slice)]
				tnode.Value = nil
				tr.length += 2
				tr.size += len(label)
				return
			}
			continue
		}
		tnode.edges = append(tnode.edges, &edge{
			label: label,
			n: &Node{
				Value:    v,
				depth:    tnode.depth + 1,
				priority: 1,
			},
		})
		tr.length++
		tr.size += len(label)
		return
	}
}

// 查找节点
func (tr *Tree) Get(label string) (*Node, map[string]string) {
	if label == "" {
		return nil, nil
	}
	if tr.safe {
		defer tr.mu.RUnlock()
		tr.mu.RLock()
	}
	tnode := tr.root
	if tr.binary {
		return tnode.getBinary(label), nil
	}
	var params map[string]string
	for tnode != nil && label != "" {
		var next *edge
	Walk:
		for _, edge := range tnode.edges {
			slice := edge.label
			for {
				phIndex := len(slice)
				// 检查是否有占位符，如果没有，则使用整个单词进行比较。
				if i := strings.IndexByte(slice, tr.placeholder); i >= 0 {
					phIndex = i
				}
				if i := strings.IndexByte(slice, '*'); i >= 0 {
					phIndex = i
				}
				prefix := slice[:phIndex]
				// 剩余部分不是占位符，继续匹配后续字符串
				if !strings.HasPrefix(label, prefix) {
					continue Walk
				}
				label = label[len(prefix):]
				// 如果"slice"是整个标签，则当前节点匹配结束，进入下一个边缘节点
				if len(prefix) == len(slice) {
					next = edge
					break Walk
				}
				// 检查是否有分隔符
				// 如果没有分隔符，将整个label进行匹配
				var delimIndex int
				slice = slice[phIndex:]
				if delimIndex = strings.IndexByte(slice[1:], tr.delim) + 1; delimIndex <= 0 {
					delimIndex = len(slice)
				}
				key := slice[1:delimIndex] // 从map键中移除占位符
				slice = slice[delimIndex:]
				if delimIndex = strings.IndexByte(label[1:], tr.delim) + 1; delimIndex <= 0 {
					delimIndex = len(label)
				}
				if params == nil {
					params = make(map[string]string)
				}
				params[key] = label[:delimIndex]
				label = label[delimIndex:]
				if slice == "" && label == "" {
					next = edge
					break Walk
				}
			}
		}
		if next != nil {
			tnode = next.n
			continue
		}
		tnode = nil
	}
	return tnode, params
}

// Del删除节点。
//如果父节点不保存任何值，则最终只保存一条边，删除一条边后，它将与剩余的边合并。
func (tr *Tree) Del(label string) {
	if string(label) == "" {
		return
	}
	if tr.safe {
		defer tr.mu.Unlock()
		tr.mu.Lock()
	}
	tnode := tr.root
	if tr.binary {
		del := tnode.delBinary(label)
		tr.length--
		bits := tr.size*8 - del
		if bits == 0 {
			tr.size = 0
			return
		}
		tr.size = (bits / 8) + 1
		return
	}
	var edgex int
	var parent *edge
	var ptrs []*int
	for tnode != nil && label != "" {
		var next *edge
		// 完全匹配查找
		for i, e := range tnode.edges {
			if strings.HasPrefix(label, e.label) {
				next = e
				edgex = i
				break
			}
		}
		if next != nil {
			tnode = next.n
			label = label[len(next.label):]
			ptrs = append(ptrs, &tnode.priority)
			if label != "" {
				parent = next
			}
			continue
		}
		// No matches.
		parent = nil
		tnode = nil
	}
	if tnode != nil {
		pnode := tr.root // 尝试根节点的标签匹配
		if parent != nil {
			pnode = parent.n
		}
		// 降低上级节点的优先级
		done := make(chan struct{})
		if tnode.Value != nil {
			go func() {
				for _, p := range ptrs {
					*p--
				}
				close(done)
			}()
		}
		// 合并tnode和父节点的边
		pnode.edges = append(pnode.edges, tnode.edges...)
		// 从父节点中删除tnode，只留下它的边缘
		pnode.edges = append(pnode.edges[:edgex], pnode.edges[edgex+1:]...)
		// 当pnode中只剩下一条边且其值为nil时，它们可以合并
		if len(pnode.edges) == 1 && pnode.Value == nil && parent != nil {
			e := pnode.edges[0]
			parent.label += e.label
			pnode.Value = e.n.Value
			pnode.edges = e.n.edges
			tr.length--
		}
		tr.length--
		if tnode.Value != nil {
			<-done
		}
	}
}

// 查询包括根的所有节点总数，
func (tr *Tree) Len() int {
	if tr.safe {
		defer tr.mu.RUnlock()
		tr.mu.RLock()
	}
	return tr.length
}

// 设置一个占位符和分隔符
func (tr *Tree) SetBoundaries(placeholder, delim byte) {
	tr.placeholder = placeholder
	tr.delim = delim
}

// 树中存储的总字节大小
func (tr *Tree) Size() int {
	return tr.size
}

// 根据它们优先级对树节点及其子节点进行递归排序
func (tr *Tree) Sort(st SortingTechnique) {
	if !tr.binary {
		if tr.safe {
			defer tr.mu.Unlock()
			tr.mu.Lock()
		}
		tr.root.sort(st)
	}
}

// 打印树结构
func (tr *Tree) String() string {
	if tr.safe {
		defer tr.mu.RUnlock()
		tr.mu.RLock()
	}
	bd := tr.bd
	bd.Reset()
	bd.WriteString(bd.colors[colorBold].Wrap("\n."))
	if tr.bd.debug {
		mag := bd.colors[colorMagenta]
		bd.WriteString(mag.Wrapf(" (%d node", tr.length))
		if tr.length != 1 {
			bd.WriteString(mag.Wrap("s")) // avoid writing "1 nodes"
		}
		bd.WriteString(mag.Wrap(")"))
	}
	tr.bd.WriteByte('\n')
	if tr.binary {
		tr.root.writeToBinary(tr.bd, &bytes.Buffer{}, &bytes.Buffer{})
	} else {
		tr.root.writeTo(tr.bd)
	}
	return tr.bd.String()
}
