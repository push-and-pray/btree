package btree

import (
	"cmp"
)

type BTree[K cmp.Ordered, V any] struct {
	degree int
	root   *Node[K, V]
}

type Node[K cmp.Ordered, V any] struct {
	children children[K, V]
	items    items[K, V]
}
type children[K cmp.Ordered, V any] []*Node[K, V]

type Item[K cmp.Ordered, V any] struct {
	key   K
	value V
}
type items[K cmp.Ordered, V any] []Item[K, V]

func (t *BTree[K, V]) minItems() int {
	return t.degree - 1
}

func (t *BTree[K, V]) maxItems() int {
	return t.degree*2 - 1
}

func (t *BTree[K, V]) maxChildren() int {
	return t.degree * 2
}

func (n *Node[K, V]) hasValidKeyChildRatio() bool {
	return len(n.items)+1 == len(n.children)
}

func (n *Node[K, V]) isLeaf() bool {
	return len(n.children) == 0
}

func NewBtree[K cmp.Ordered, V any](degree int) *BTree[K, V] {
	if degree < 2 {
		panic("Invalid degree. Must be larger than 1")
	}
	bt := BTree[K, V]{degree: degree}
	return &bt
}

func (t *BTree[K, V]) newNode() *Node[K, V] {
	return &Node[K, V]{
		children: make([]*Node[K, V], 0, t.maxChildren()),
		items:    make([]Item[K, V], 0, t.maxItems()),
	}
}

/*
Attempt to get item with key k. Success is indicated by returned bool
*/
func (t *BTree[K, V]) Get(k K) (V, bool) {
	if t.root == nil {
		var zeroVal V
		return zeroVal, false
	}
	item, found := t.get(k, t.root)
	return item.value, found
}

/*
Attempt to get item with key k from subtree rooted at n. Success is indicated by returned bool
*/
func (t *BTree[K, V]) get(k K, n *Node[K, V]) (Item[K, V], bool) {
	idx, found := n.items.find(k)

	if found {
		return n.items[idx], true
	}

	if n.isLeaf() {
		var zeroVal Item[K, V]
		return zeroVal, false
	}

	return t.get(k, n.children[idx])

}

/*
Splits a node n. Returns the promoted item and the new node
*/
func (t *BTree[K, V]) split(n *Node[K, V]) (Item[K, V], *Node[K, V]) {
	median := len(n.items) / 2

	promotedItem := n.items[median]
	newNode := t.newNode()
	newNode.items = append(newNode.items, n.items[median+1:]...)
	n.items = n.items[:median]

	if !n.isLeaf() {
		newNode.children = append(newNode.children, n.children[median+1:]...)
		n.children = n.children[:median+1]
	}

	return promotedItem, newNode
}

/*
Insert key,value pair into btree
*/
func (t *BTree[K, V]) Insert(k K, v V) {
	// Initialize btree if required
	if t.root == nil {
		t.root = t.newNode()
		t.root.items = append(t.root.items, Item[K, V]{k, v})
		return
	}
	if len(t.root.items) >= t.maxItems() {
		promotedItem, splitNode := t.split(t.root)
		newRoot := t.newNode()
		newRoot.items = append(newRoot.items, promotedItem)
		newRoot.children = append(newRoot.children, t.root, splitNode)
		t.root = newRoot
	}

	t.insert(k, v, t.root)
}

/*
Insert key, value pair into subtree rooted at n. Returns information regarding if
the insertion resulted in a promotion, which the caller must handle
*/
func (t *BTree[K, V]) insert(k K, v V, n *Node[K, V]) {
	idx, found := n.items.find(k)

	// If the key already exists, replace it
	if found {
		n.items[idx].value = v
		return
	}

	if n.isLeaf() {
		n.items.insertAt(k, v, idx)
		return
	}

	next := n.children[idx]
	if len(next.items) >= t.maxItems() {
		promotedItem, splitNode := t.split(next)
		n.items.insertAt(promotedItem.key, promotedItem.value, idx)
		n.children.insertAt(splitNode, idx+1)

		idx, found = n.items.find(k)
		if found {
			n.items[idx].value = v
			return
		}

	}

	t.insert(k, v, n.children[idx])
}

/*
Delete item with key k from btree. Returns whether the key was found
*/
func (t *BTree[K, V]) Delete(k K) bool {

	if t.root == nil {
		return false
	}

	found := t.delete(k, t.root)
	if !found {
		return false
	}

	// Handle shrinking of btree
	if len(t.root.items) == 0 {
		if t.root.isLeaf() {
			t.root = nil
		} else {
			t.root = t.root.children[0]
		}
	}

	return true
}

/*
Delete item with key k from subtree rooted at n. Returns whether key was found
*/
func (t *BTree[K, V]) delete(k K, n *Node[K, V]) bool {
	idx, found := n.items.find(k)
	if found {
		if n.isLeaf() {
			n.items.deleteAt(idx)
			return true
		}

		// Deletion from internal nodes are split into three cases

		// If the left child has enough items, we will replace the
		// key which its predecessor. The same can be possible for
		// The right child. If neither of them have enough, we must merge
		if leftChild := n.children[idx]; len(leftChild.items) > t.minItems() {
			if len(leftChild.items) <= t.minItems() {
				leftChild = t.rebalance(n, idx)
			}
			n.items[idx] = t.popMax(leftChild)
		} else if rightChild := n.children[idx+1]; len(rightChild.items) > t.minItems() {
			if len(rightChild.items) <= t.minItems() {
				rightChild = t.rebalance(n, idx+1)
			}
			n.items[idx] = t.popMin(rightChild)
		} else {
			n.merge(idx)
			t.delete(k, leftChild)

		}

		return true
	}

	// If we are at a leaf, and we still havent found the key, it is not here
	if n.isLeaf() {
		return false
	}

	// Recurse further, ensuring that every child we recurse into
	// has more than minimum amount of items
	child := n.children[idx]
	if len(child.items) > t.minItems() {
		return t.delete(k, child)
	}

	child = t.rebalance(n, idx)

	return t.delete(k, child)

}

/*
Rebalances child at index i of node n. Returns a pointer to child i
or its left sibling, if child i got merged into it
*/
func (t *BTree[K, V]) rebalance(n *Node[K, V], i int) *Node[K, V] {

	hasLeftSibling := i > 0
	hasRightSibling := i < len(n.children)-1

	if hasLeftSibling && len(n.children[i-1].items) > t.minItems() {
		n.stealFromLeftSibling(i)
	} else if hasRightSibling && len(n.children[i+1].items) > t.minItems() {
		n.stealFromRightSibling(i)
	} else {
		if hasRightSibling {
			n.merge(i)
		} else {
			n.merge(i - 1)
			// We have merged our old target into its left sibling and must change course
			return n.children[i-1]
		}

	}
	return n.children[i]
}

/*
Pop the max item at the btree rooted at node n, assuming that n has more than min items
*/
func (t *BTree[K, V]) popMax(n *Node[K, V]) Item[K, V] {
	if n.isLeaf() {
		return n.items.deleteAt(len(n.items) - 1)
	}

	next := n.children[len(n.children)-1]
	if len(next.items) <= t.minItems() {
		next = t.rebalance(n, len(n.children)-1)
	}
	return t.popMax(next)
}

/*
Pop the min item at the btree rooted at node n, assuming that n has more than min items
*/
func (t *BTree[K, V]) popMin(n *Node[K, V]) Item[K, V] {
	if n.isLeaf() {
		return n.items.deleteAt(0)
	}

	next := n.children[0]

	if len(next.items) <= t.minItems() {
		next = t.rebalance(n, 0)
	}
	return t.popMin(next)
}

// Steals an item from the left sibling of child at index i of node n
func (n *Node[K, V]) stealFromLeftSibling(i int) {
	child, sibling := n.children[i], n.children[i-1]
	demotedItem := n.items[i-1]
	child.items.insertAt(demotedItem.key, demotedItem.value, 0)
	if !sibling.isLeaf() {

		siblingChild := sibling.children.deleteAt(len(sibling.children) - 1)

		child.children.insertAt(siblingChild, 0)
	}
	promotedItem := sibling.items.deleteAt(len(sibling.items) - 1)
	n.items[i-1] = promotedItem
}

// Steals an item from the right sibling of child at index i of node n
func (n *Node[K, V]) stealFromRightSibling(i int) {
	child, sibling := n.children[i], n.children[i+1]
	child.items = append(child.items, n.items[i])
	if !child.isLeaf() {
		child.children = append(child.children, sibling.children.deleteAt(0))
	}
	n.items[i] = sibling.items.deleteAt(0)

}

// Merge child at index i of node n, with child at index i+1
func (n *Node[K, V]) merge(i int) {
	child, sibling := n.children[i], n.children[i+1]

	child.items = append(child.items, n.items.deleteAt(i))
	child.items = append(child.items, sibling.items...)

	if !child.isLeaf() {
		child.children = append(child.children, sibling.children...)
	}
	n.children.deleteAt(i + 1)

}
