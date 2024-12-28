package btree

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

func (btree *BTree[K, V]) checkTreeValid(node *Node[K, V], t *testing.T) {
	isLeaf := len(node.children) == 0
	isRoot := btree.root == node

	if len(node.children) > btree.maxChildren() {
		t.Errorf("Node has too many children: %+v", *node)
	}

	if len(node.items) < btree.minItems() && !isRoot {
		t.Errorf("Node has too few items: %+v", *node)
	}

	if len(node.items) > btree.maxItems() {
		t.Errorf("Node has too many keys: %+v", *node)
	}

	isItemsSorted := slices.IsSortedFunc(node.items, func(a Item[K, V], b Item[K, V]) int {
		return cmp.Compare(a.key, b.key)
	})
	if !isItemsSorted {
		t.Errorf("Items of node are not sorted: %+v", *node)
	}

	if !node.hasValidKeyChildRatio() && !isLeaf && !isRoot {
		t.Errorf("Node doesn't have valid number of children vs items: %+v", *node)
	}

	for _, child := range node.children {
		btree.checkTreeValid(child, t)
	}

}

func (btree *BTree[K, V]) _hasValidDepth(node *Node[K, V], depth int, leafDepth *int) bool {
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

func (btree *BTree[K, V]) hasValidDepth(t *testing.T) {
	leafDepth := -1
	if !btree._hasValidDepth(btree.root, 0, &leafDepth) {
		t.Errorf("BTree does not have a valid depth")
	}
}

func newPopulatedBTree() *BTree[int, string] {

	layer_3_1 := &Node[int, string]{
		items: []Item[int, string]{{59, "59"}, {61, "61"}},
	}
	layer_3_2 := &Node[int, string]{
		items: []Item[int, string]{{71, "71"}, {79, "79"}, {83, "83"}},
	}
	layer_3_3 := &Node[int, string]{
		items: []Item[int, string]{{97, "97"}, {101, "101"}},
	}

	layer_2_1 := &Node[int, string]{
		items: []Item[int, string]{{5, "5"}, {7, "7"}},
	}
	layer_2_2 := &Node[int, string]{
		items: []Item[int, string]{{17, "17"}, {23, "23"}},
	}
	layer_2_3 := &Node[int, string]{
		items: []Item[int, string]{{31, "31"}, {37, "37"}},
	}

	layer_1_0 := &Node[int, string]{
		items:    []Item[int, string]{{11, "11"}, {29, "29"}},
		children: []*Node[int, string]{layer_2_1, layer_2_2, layer_2_3},
	}
	layer_1_1 := &Node[int, string]{
		items:    []Item[int, string]{{67, "67"}, {89, "89"}},
		children: []*Node[int, string]{layer_3_1, layer_3_2, layer_3_3},
	}

	root := &Node[int, string]{
		items:    []Item[int, string]{{43, "43"}},
		children: []*Node[int, string]{layer_1_0, layer_1_1},
	}

	return &BTree[int, string]{degree: 3, root: root}

}

func TestBTreeGet(t *testing.T) {
	btree := newPopulatedBTree()

	tests := []struct {
		key              int
		expected_val     string
		expected_success bool
	}{
		{5, "5", true},
		{7, "7", true},
		{17, "17", true},
		{23, "23", true},
		{31, "31", true},
		{37, "37", true},
		{43, "43", true},
		{59, "59", true},
		{61, "61", true},
		{67, "67", true},
		{71, "71", true},
		{79, "79", true},
		{83, "83", true},
		{89, "89", true},
		{97, "97", true},
		{101, "101", true},
		{696969, "", false},
		{1, "", false},
		{72, "", false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Get(%d)", test.key), func(t *testing.T) {
			result, success := btree.Get(test.key)
			if result != test.expected_val || success != test.expected_success {
				t.Errorf("expected %s and %t, got %s and %t", test.expected_val, test.expected_success, result, success)
			}
			btree.checkTreeValid(btree.root, t)
			btree.hasValidDepth(t)
		})
	}
}

func TestBTreeInsert(t *testing.T) {
	r := rand.NewPCG(424242, 1024)
	random := rand.New(r)

	randomInserts := 1000

	type testInput struct {
		key   int
		value string
	}

	randomInputs := make([]testInput, 0, randomInserts)

	for range randomInserts {
		key := random.IntN(1000)
		test1 := testInput{key, strconv.Itoa(key)}
		randomInputs = append(randomInputs, test1)
	}

	btree := NewBtree[int, string](3)
	t.Run("Random inserts", func(t *testing.T) {

		for _, test := range randomInputs {
			btree.Add(test.key, test.value)
			btree.checkTreeValid(btree.root, t)
			btree.hasValidDepth(t)
		}
	})
}

func BenchmarkBTreeGetLayer1(b *testing.B) {
	btree := newPopulatedBTree()

	keys := []int{43}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Get(key)
		}
	}
}

func BenchmarkBTreeGetLayer2(b *testing.B) {
	btree := newPopulatedBTree()

	keys := []int{11, 29, 67, 89}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Get(key)
		}
	}
}

func BenchmarkBTreeGetLayer3(b *testing.B) {
	btree := newPopulatedBTree()

	keys := []int{5, 7, 17, 23, 31, 37, 59, 61, 71, 79, 83, 97, 101}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Get(key)
		}
	}
}

func BenchmarkBTreeGetNonExisting(b *testing.B) {
	btree := newPopulatedBTree()

	keys := []int{696969, 1, 72}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Get(key)
		}
	}
}
