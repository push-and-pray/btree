package btree

import (
	"fmt"
	"sort"
	"strings"
)

/*
Returns the index where a key should be inserted in a items slice. Returns a true bool if the same key was found
*/

func (s items[K, V]) find(k K) (int, bool) {
	idx := sort.Search(len(s), func(i int) bool {
		return s[i].key >= k
	})

	if idx < len(s) && s[idx].key == k {
		return idx, true
	}

	return idx, false
}

/*
insertAt inserts a new item with key k and value v at the specified index i
in the slice of items. It shifts the elements at and after index i to the
right to make space for the new item.
*/
func (s *items[K, V]) insertAt(k K, v V, i int) {
	var zeroVal Item[K, V]
	*s = append(*s, zeroVal)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = Item[K, V]{k, v}
}

func (s *children[K, V]) insertAt(node *Node[K, V], i int) {
	*s = append(*s, nil)
	copy((*s)[i+1:], (*s)[i:])
	(*s)[i] = node
}

func (btree *BTree[K, V]) String() string {
	var sb strings.Builder
	btree.stringHelper(btree.root, 0, &sb)
	return sb.String()
}

func (btree *BTree[K, V]) stringHelper(node *Node[K, V], level int, sb *strings.Builder) {
	if node == nil {
		return
	}

	indent := strings.Repeat("  ", level)
	sb.WriteString(indent + "Node: ")

	for _, item := range node.items {
		sb.WriteString(fmt.Sprintf("%v:%v ", item.key, item.value))
	}
	sb.WriteString("\n")

	for _, child := range node.children {
		btree.stringHelper(child, level+1, sb)
	}
}
