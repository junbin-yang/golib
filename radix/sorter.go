package radix

type SortingTechnique uint8

const (
	// 树的边缘按标签升序排列
	AscLabelSort SortingTechnique = iota
	// 树的边缘按标签降序排列
	DescLabelSort
	// 树的边缘按优先级排序
	PrioritySort
)

type sorter struct {
	n  *Node
	st SortingTechnique
}

func (s *sorter) Len() int {
	return len(s.n.edges)
}

func (s *sorter) Less(i, j int) bool {
	n := s.n
	switch s.st {
	case AscLabelSort:
		return n.edges[i].label < n.edges[j].label
	case DescLabelSort:
		return n.edges[i].label > n.edges[j].label
	default:
		return n.edges[i].n != nil &&
			n.edges[j].n != nil &&
			n.edges[i].n.priority > n.edges[j].n.priority
	}
}

func (s *sorter) Swap(i, j int) {
	s.n.edges[i], s.n.edges[j] = s.n.edges[j], s.n.edges[i]
}
