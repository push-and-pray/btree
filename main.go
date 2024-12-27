package main

import (
	"cmp"
	"fmt"
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
	node := btree.root
	assert(func() bool { return node != nil }, "Empty BTree encountered!")

nodeTraversalLoop:
	for {
		for idx, elm := range node.elements {
			switch cmp.Compare(k, elm.key) {
			case 0:
				return elm.value, true
			case -1:
				assert(func() bool { return len(node.children) > idx && idx >= 0 }, "Attempted illegal traversel")
				node = node.children[idx]
				continue nodeTraversalLoop
			}
		}
		if len(node.children) == 0 {
			var ZeroVal V
			return ZeroVal, false
		}
		assert(func() bool { return len(node.children) > len(node.elements) && len(node.elements) >= 0 }, "Attempted illegal traversel")
		node = node.children[len(node.elements)]
	}
}

func (btree *BTree[K, V]) Add(k K, v V) {
	node := btree.root
	assert(func() bool { return node != nil }, "Empty BTree encountered!")

	// Find correct leaf to insert key at
nodeTraversalLoop:
	for {
		if len(node.children) == 0 {
			break nodeTraversalLoop
		}

		for idx, elm := range node.elements {
			switch cmp.Compare(k, elm.key) {
			case 0:
				// Return early if we find a matching key
				node.elements[idx].value = v
				return
			case -1:
				assert(func() bool { return len(node.children) > idx && idx >= 0 }, "Attempted illegal traversel")
				node = node.children[idx]
				continue nodeTraversalLoop
			}
		}

		assert(func() bool { return len(node.children) > len(node.elements) && len(node.elements) >= 0 }, "Attempted illegal traversel")
		node = node.children[len(node.elements)]
	}

	// If there is space in the node, insert it, while keeping order
	if len(node.elements) < btree.max {
		node.elements, _ = btree.insertOrdered(node.elements, k, v)
		return
	}

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
