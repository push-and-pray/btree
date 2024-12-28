package btree

import (
	"cmp"
)

type BTree[K cmp.Ordered, V any] struct {
	degree int // 2
	root   *Node[K, V]
}

func (btree *BTree[K, V]) minItems() int {
	return btree.degree - 1 // 1
}

func (btree *BTree[K, V]) maxItems() int {
	return btree.degree*2 - 1 // 3
}

func (btree *BTree[K, V]) maxChildren() int {
	return btree.maxItems() + 1 // 4
}

type children[K cmp.Ordered, V any] []*Node[K, V]
type Node[K cmp.Ordered, V any] struct {
	children children[K, V]
	items    items[K, V]
}

func (node *Node[K, V]) hasValidKeyChildRatio() bool {
	return len(node.items)+1 == len(node.children)
}

type items[K cmp.Ordered, V any] []Item[K, V]
type Item[K cmp.Ordered, V any] struct {
	key   K
	value V
}

func NewBtree[K cmp.Ordered, V any](degree int) *BTree[K, V] {
	bt := BTree[K, V]{degree: degree}
	return &bt
}

func (btree *BTree[K, V]) newNode() *Node[K, V] {
	return &Node[K, V]{
		children: make([]*Node[K, V], 0, btree.maxChildren()),
		items:    make([]Item[K, V], 0, btree.maxItems()),
	}
}

/*
Splits Node n at item index i. Return the split item and the new node
*/

func (btree *BTree[K, V]) Get(k K) (V, bool) {
	return btree.get(k, btree.root)
}

func (btree *BTree[K, V]) get(k K, root *Node[K, V]) (V, bool) {
	idx, found := root.items.find(k)

	if found {
		return root.items[idx].value, true
	}

	if len(root.children) == 0 {
		var zeroVal V
		return zeroVal, false
	}

	return btree.get(k, root.children[idx])

}

func (btree *BTree[K, V]) split(n *Node[K, V]) (Item[K, V], *Node[K, V]) {
	if len(n.items) < btree.maxItems() {
		panic("Tried to split non full node")
	}
	median := len(n.items) / 2

	splitItem := n.items[median]
	newNode := btree.newNode()
	newNode.items = append(newNode.items, n.items[median+1:]...)
	n.items = n.items[:median]

	if len(n.children) > 0 {
		newNode.children = append(newNode.children, n.children[median+1:]...)
		n.children = n.children[:median+1]
	}

	return splitItem, newNode
}

func (btree *BTree[K, V]) add(k K, v V, n *Node[K, V]) {
	idx, found := n.items.find(k)

	if found {
		n.items[idx].value = v
		return
	}

	if len(n.children) == 0 {
		n.items.insertAt(k, v, idx)
		return
	}

	// Split the child which we are traversing towards, if necessary
	if len(n.children[idx].items) >= btree.maxItems() {
		medianItem, newNode := btree.split(n.children[idx])
		n.items.insertAt(medianItem.key, medianItem.value, idx)
		n.children.insertAt(newNode, idx+1)

		inTree := n.items[idx]

		switch cmp.Compare(k, inTree.key) {
		case -1:
			break
		case 1:
			idx++
		case 0:
			n.items[idx].value = v
			return
		}

	}
	btree.add(k, v, n.children[idx])
}

func (btree *BTree[K, V]) Add(k K, v V) {

	if btree.root == nil {
		btree.root = btree.newNode()
		btree.root.items = append(btree.root.items, Item[K, V]{k, v})
		return
	}

	if len(btree.root.items) >= btree.maxItems() {
		medianItem, newNode := btree.split(btree.root)
		oldRoot := btree.root
		btree.root = btree.newNode()
		btree.root.items = append(btree.root.items, medianItem)
		btree.root.children = append(btree.root.children, oldRoot, newNode)
	}

	btree.add(k, v, btree.root)

}
