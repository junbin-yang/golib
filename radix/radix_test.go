package radix

import (
	"fmt"
	"testing"
)

func TestBinaryTree(t *testing.T) {
	tree := BinaryRoot()

	tree.Insert("abc", 1)
	tree.Insert("abdd", 2)
	tree.Insert("app", 3)
	tree.Insert("ppd", 4)
	tree.Insert("abcdefg", 5)
	tree.Insert("ccd", 6)
	tree.Insert("aew", 7)
	tree.Sort(PrioritySort)
	n, _ := tree.Search("app")
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
		tree.Insert(route, i)
	}

	requests := []string{
		"/index/zhangsan",
		"/path/456/user/foobar",
	}

	for _, request := range requests {
		n, v := tree.Search(request)
		t.Log(fmt.Printf("%+v ,%+v\n", n, v))
	}

}
