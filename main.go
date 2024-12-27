package main

import (
	"cmp"
	"fmt"
	"sort"
	"strings"
)

type BTree[K cmp.Ordered, V any] struct {
	max  int // Max number of elements
	root *Node[K, V]
}

type Node[K cmp.Ordered, V any] struct {
	children []*Node[K, V]
	elements []Element[K, V]
}

type Element[K cmp.Ordered, V any] struct {
	key   K
	value V
}

func NewBtree[K cmp.Ordered, V any](max int) BTree[K, V] {
	return BTree[K, V]{
		max: max,
		root: &Node[K, V]{
			children: make([]*Node[K, V], 0, max+1),
			elements: make([]Element[K, V], 0, max),
		},
	}
}

func (btree *BTree[K, V]) newNode() Node[K, V] {
	return Node[K, V]{
		children: make([]*Node[K, V], 0, btree.max+1),
		elements: make([]Element[K, V], 0, btree.max),
	}
}

func (btree *BTree[K, V]) Get(k K) (V, bool) {
	return btree.get(k, btree.root)
}

func (btree *BTree[K, V]) get(k K, root *Node[K, V]) (V, bool) {
	idx := sort.Search(len(root.elements), func(i int) bool { return root.elements[i].key >= k })

	if idx < len(root.elements) && root.elements[idx].key == k {
		return root.elements[idx].value, true
	}

	if len(root.children) == 0 {
		var zeroVal V
		return zeroVal, false
	}

	return btree.get(k, root.children[idx])

}

func (btree *BTree[K, V]) Add(k K, v V) {
	node := btree.root
	assert(func() bool { return node != nil },
		"Empty BTree root encountered!")

	if len(node.elements) > btree.max {
		// newRoot := Node[K, V]{children: []*Node[K, V]{node}}

	}
}
func (btree *BTree[K, V]) splitChild(parent *Node[K, V], index int) {
	//child := parent.children[index]
	//newNode := &Node[K, V]{}

}

func (btree *BTree[K, V]) String() string {
	if btree.root == nil {
		return "Empty B-Tree"
	}
	return nodeToString(btree.root, 0)
}

func nodeToString[K cmp.Ordered, V any](node *Node[K, V], level int) string {
	var sb strings.Builder

	indent := strings.Repeat("  ", level)
	sb.WriteString(fmt.Sprintf("%sNode (level %d):\n", indent, level))
	for _, elem := range node.elements {
		sb.WriteString(fmt.Sprintf("%s  Key: %v, Value: %v\n", indent, elem.key, elem.value))
	}

	for _, child := range node.children {
		sb.WriteString(nodeToString(child, level+1))
	}
	return sb.String()
}

func main() {
	btree := NewBtree[int, string](4)
	print(btree.String())
}
