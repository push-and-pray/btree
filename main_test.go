package btree

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

func (btree *BTree[K, V]) checkTreeValid(node *Node[K, V], t *testing.T) bool {
	if node == nil {
		return true
	}
	isLeaf := len(node.children) == 0
	isRoot := btree.root == node

	valid := true

	if len(node.children) > btree.maxChildren() {
		t.Errorf("Node has too many children: %+v", *node)
		valid = false
	}

	if len(node.items) < btree.minItems() && !isRoot {
		t.Errorf("Node has too few items: %+v", *node)
		valid = false
	}

	if len(node.items) > btree.maxItems() {
		t.Errorf("Node has too many keys: %+v", *node)
		valid = false
	}

	isItemsSorted := slices.IsSortedFunc(node.items, func(a Item[K, V], b Item[K, V]) int {
		return cmp.Compare(a.key, b.key)
	})
	if !isItemsSorted {
		t.Errorf("Items of node are not sorted: %+v", *node)
		valid = false
	}

	if !node.hasValidKeyChildRatio() && !isLeaf && !isRoot {
		t.Errorf("Node doesn't have valid number of children vs items: %+v", *node)
		valid = false
	}

	if !node.isLeaf() {
		for idx, item := range node.items {
			for _, predecessor := range node.children[idx].items {
				if predecessor.key >= item.key {
					t.Errorf("Predecessor %v is larger than item %v", predecessor.key, item.key)
					valid = false

				}
			}

			for _, successor := range node.children[idx+1].items {
				if successor.key <= item.key {
					t.Errorf("Successor %v is smaller than item %v", successor.key, item.key)
					valid = false
				}
			}
		}
	}

	for _, child := range node.children {
		if !btree.checkTreeValid(child, t) {
			valid = false
		}
	}

	return valid
}

func (btree *BTree[K, V]) _hasValidDepth(node *Node[K, V], depth int, leafDepth *int) bool {
	if btree.root == nil {
		return true
	}
	if len(node.children) == 0 {
		if *leafDepth == -1 {
			*leafDepth = depth
		}
		return depth == *leafDepth
	}

	for _, child := range node.children {
		if !btree._hasValidDepth(child, depth+1, leafDepth) {
			return false
		}
	}
	return true
}

func (btree *BTree[K, V]) hasValidDepth(t *testing.T) bool {
	leafDepth := -1
	if !btree._hasValidDepth(btree.root, 0, &leafDepth) {
		t.Errorf("BTree does not have a valid depth")
		return false
	}
	return true
}

func TestBTreeGet(t *testing.T) {
	for d := 2; d < 10; d++ {
		btree := NewBtree[int, string](d)

		for i := range 50 {
			if i == 25 {
				continue
			}
			btree.Insert(i, strconv.Itoa(i))
		}

		tests := []struct {
			key      int
			expected string
			found    bool
		}{
			{0, "0", true},
			{49, "49", true},
			{22, "22", true},
			{-1, "", false},
			{50, "", false},
			{25, "", false},
		}

		t.Run(fmt.Sprintf("Get at degree %v", d), func(t *testing.T) {
			for _, test := range tests {
				val, found := btree.Get(test.key)
				if found != test.found {
					t.Errorf("Get(%d) found = %v; want %v", test.key, found, test.found)
				}
				if found && val != test.expected {
					t.Errorf("Get(%d) = %v; want %v", test.key, val, test.expected)
				}
				btree.checkTreeValid(btree.root, t)
				btree.hasValidDepth(t)

			}
		})
	}
}

func TestBTreeRandomInserts(t *testing.T) {
	r := rand.NewPCG(424242, 1024)
	random := rand.New(r)

	type testInput struct {
		key   int
		value string
	}

	for i := 2; i < 10; i++ {
		randomInserts := int(math.Pow(float64(i), float64(4)))

		for run := range 100 {
			randomInputs := make([]testInput, 0, randomInserts)
			for range randomInserts {
				key := random.IntN(100)
				value := random.IntN(100)
				test1 := testInput{key, strconv.Itoa(value)}
				randomInputs = append(randomInputs, test1)
			}

			btree := NewBtree[int, string](i)
			t.Run(fmt.Sprintf("Random inserts at degree %v #%v", i, run), func(t *testing.T) {
				for _, test := range randomInputs {
					btree.Insert(test.key, test.value)
					t.Logf("State:\n%v", btree.String())

					treeValid := btree.checkTreeValid(btree.root, t)
					validDepth := btree.hasValidDepth(t)
					fail := false
					if !treeValid || !validDepth {
						fail = true
						t.Error("Tree is not valid!")
					}

					val, found := btree.Get(test.key)
					if !found {
						fail = true
						t.Errorf("Key disappeared, expected: %v, got nothing", test.value)
					}
					if val != test.value {
						fail = true
						t.Errorf("Unexpected value returned: expected: %v, got: %v", test.value, val)
					}

					if fail {
						t.Fatal()
					}
				}
			})
		}
	}
}

func TestBTreeRandomDeletes(t *testing.T) {
	r := rand.NewPCG(424242, 1024)
	random := rand.New(r)

	type testInput struct {
		key   int
		value string
	}

	for i := 2; i < 10; i++ {

		randomDeletes := int(math.Pow(float64(i), float64(4)))

		for run := range 100 {
			randomInputs := make([]testInput, 0, randomDeletes)
			for range randomDeletes {
				key := random.IntN(100)
				value := random.IntN(100)
				test1 := testInput{key, strconv.Itoa(value)}
				randomInputs = append(randomInputs, test1)
			}
			btree := NewBtree[int, string](i)
			t.Run(fmt.Sprintf("Random deletes at degree %v #%v", i, run), func(t *testing.T) {

				for _, test := range randomInputs {
					btree.Insert(test.key, test.value)
				}
				t.Log("Starting from this tree:")
				t.Logf("\n%v", btree.String())

				found := btree.Delete(-1)
				if found {
					t.Fatalf("Deleted non-existing key with")
				}

				for _, test := range randomInputs {
					t.Logf("Deleting key: %v", test.key)
					btree.Delete(test.key)

					_, found := btree.Get(test.key)

					treeValid := btree.checkTreeValid(btree.root, t)
					validDepth := btree.hasValidDepth(t)
					fail := false
					if !treeValid || !validDepth {
						fail = true
						t.Error("Tree is not valid!")
					}

					if found {
						fail = true
						t.Errorf("Key %v still here %v", test.key, run)
					}

					if fail {
						t.Logf("\n%v", btree.String())
						t.Fatal()
					}
				}

				found = btree.Delete(-1)
				if found {
					t.Fatalf("Deleted non-existing key")
				}
			})
		}
	}
}
