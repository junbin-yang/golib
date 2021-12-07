package radix

import (
	"fmt"
	"testing"
)

func TestBinaryTree(t *testing.T) {
	tree := BinaryRoot()

	tree.Add("abc", 1)
	tree.Add("abdd", 2)
	tree.Add("app", 3)
	tree.Add("ppd", 4)
	tree.Add("abcdefg", 5)
	tree.Add("ccd", 6)
	tree.Add("aew", 7)
	tree.Sort(PrioritySort)
	n, _ := tree.Get("app")
	t.Log(tree, n.Value)
}

func TestPrefixTree(t *testing.T) {
	tree := PrefixRoot()
	routes := [...]string{
		"/foozip",
		"/index/:name/ha",
		"/path/:id/index/:name",
		"/path/:id/user/:name",
		"/path/:id/user/*",
	}
	for i, route := range routes {
		tree.Add(route, i)
	}

	requests := []string{
		"/index/zhangsan",
		"/path/456/user/foobar",
	}

	for _, request := range requests {
		n, v := tree.Get(request)
		t.Log(fmt.Printf("%+v ,%+v\n", n, v))
	}

}
