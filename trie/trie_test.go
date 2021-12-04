package trie

import (
	"fmt"
	"sort"
	"testing"
)

func TestTrieAdd(t *testing.T) {
	trie := NewTrie()

	n := trie.Add("POST", 1)

	if n.Meta().(int) != 1 {
		t.Errorf("Expected 1, got: %d", n.Meta().(int))
	}
}

func TestTrieFind(t *testing.T) {
	trie := NewTrie()
	trie.Add("POST", 1)

	n, err := trie.Find("POST")
	if err != nil {
		t.Fatal("Could not find node", err.Error())
	}

	if n.Meta().(int) != 1 {
		t.Errorf("Expected 1, got: %d", n.Meta().(int))
	}
}

func TestRemoveAndFindKeys(t *testing.T) {
	trie := NewTrie()
	initial := []string{"football", "foostar", "foosball"}

	for _, key := range initial {
		trie.Add(key, nil)
	}

	trie.Remove("foosball")
	keys := trie.Keys()

	for _, k := range keys {
		if k != "football" && k != "foostar" {
			t.Errorf("key was: %s", k)
		}
	}

	keys = trie.FuzzySearch("footb")
	for i, k := range keys {
		if k != "football" && k != "foostar" {
			t.Errorf("Expected football got: %#v", k)
		}
		t.Log(k, initial[i])
	}
}

// 前缀搜索
func TestPrefixSearch(t *testing.T) {
	trie := NewTrie()
	expected := []string{
		"foo",
		"foosball",
		"football",
		"foreboding",
		"forementioned",
		"foretold",
		"foreverandeverandeverandever",
		"forbidden",
	}

	trie.Add("bar", nil)
	for _, key := range expected {
		trie.Add(key, nil)
	}

	tests := []struct {
		pre      string
		expected []string
		length   int
	}{
		{"fo", expected, len(expected)},
		{"foosbal", []string{"foosball"}, 1},
		{"abc", []string{}, 0},
	}

	for _, test := range tests {
		actual := trie.PrefixSearch(test.pre)
		sort.Strings(actual)
		sort.Strings(test.expected)

		if len(actual) > 0 {
			for i, key := range actual {
				if key != test.expected[i] {
					t.Errorf("Expected %v got: %v", test.expected[i], key)
				}
				t.Log(key)
			}
			fmt.Println()
		} else {
			fmt.Println("PrefixSearch", test.pre, "is nil")
		}
	}
}

// 模糊搜素
func TestFuzzySearch(t *testing.T) {
	trie := NewTrie()
	setup := []string{
		"foosball",
		"football",
		"bmerica",
		"ked",
		"kedlock",
		"frosty",
		"bfrza",
		"foo/bart/baz.go",
	}

	for _, key := range setup {
		trie.Add(key, nil)
	}

	tests := []string{"fsb", "footbal", "football", "fs", "oos", "kl", "ft", "fy", "fz", "a"}

	for _, test := range tests {
		actual := trie.FuzzySearch(test)
		t.Log(actual)
	}
}

func TestFuzzySearchSorting(t *testing.T) {
	trie := NewTrie()
	setup := []string{
		"foosball",
		"football",
		"bmerica",
		"ked",
		"kedlock",
		"frosty",
		"bfrza",
		"foo/bart/baz.go",
	}

	for _, key := range setup {
		trie.Add(key, nil)
	}

	actual := trie.FuzzySearch("fz")
	expected := []string{"bfrza", "foo/bart/baz.go"}

	for i, v := range expected {
		if actual[i] != v {
			t.Errorf("Expected %s got %s", v, actual[i])
		}
		t.Log(v)
	}
}

/*
func TestRoute(t *testing.T) {
	trie := NewTrie()
	// /index/:user
	trie.Add("/index/:user", 12)
	trie.Add("/:index/:user", 13)
	trie.Add("/index/:user/:id/*p", "uid")
	trie.Add("/index/:user/:name/:id", "uid")

	actual := trie.FuzzySearch("/::*")

	fmt.Printf("%+v\n", actual)
	//fmt.Println(trie.Find(actual[0]))
}
*/
