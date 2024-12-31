package btree

import (
	"cmp"
)

type BTree[K cmp.Ordered, V any] struct {
	degree int
	root   *Node[K, V]
}

func (btree *BTree[K, V]) minItems() int {
	return btree.degree - 1
}

func (btree *BTree[K, V]) maxItems() int {
	return btree.degree*2 - 1
}

func (btree *BTree[K, V]) maxChildren() int {
	return btree.degree * 2
}

type children[K cmp.Ordered, V any] []*Node[K, V]
type Node[K cmp.Ordered, V any] struct {
	children children[K, V]
	items    items[K, V]
}

func (node *Node[K, V]) hasValidKeyChildRatio() bool {
	return len(node.items)+1 == len(node.children)
}

func (node *Node[K, V]) isLeaf() bool {
	return len(node.children) == 0
}

type items[K cmp.Ordered, V any] []Item[K, V]
type Item[K cmp.Ordered, V any] struct {
	key   K
	value V
}

func NewBtree[K cmp.Ordered, V any](degree int) *BTree[K, V] {
	if degree < 2 {
		panic("Invalid degree. Must be larger than 1")
	}
	bt := BTree[K, V]{degree: degree}
	return &bt
}

func (btree *BTree[K, V]) newNode() *Node[K, V] {
	return &Node[K, V]{
		children: make([]*Node[K, V], 0, btree.maxChildren()),
		items:    make([]Item[K, V], 0, btree.maxItems()),
	}
}

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

func (btree *BTree[K, V]) Insert(k K, v V) {
	if btree.root == nil {
		btree.root = btree.newNode()
		btree.root.items = append(btree.root.items, Item[K, V]{k, v})
		return
	}

	promotedItem, splitNode, promoted := btree.insert(k, v, btree.root)
	if !promoted {
		return
	}

	newRoot := btree.newNode()
	newRoot.items = append(newRoot.items, promotedItem)
	newRoot.children = append(newRoot.children, btree.root, splitNode)
	btree.root = newRoot

}

func (btree *BTree[K, V]) insert(k K, v V, node *Node[K, V]) (Item[K, V], *Node[K, V], bool) {
	idx, found := node.items.find(k)
	var zeroVal Item[K, V]

	if found {
		node.items[idx].value = v
		return zeroVal, nil, false
	}

	var promotedItem Item[K, V]
	var splitNode *Node[K, V]
	var promoted bool
	if len(node.children) != 0 {
		promotedItem, splitNode, promoted = btree.insert(k, v, node.children[idx])
		if !promoted {
			return promotedItem, splitNode, promoted
		}
	} else {
		node.items.insertAt(k, v, idx)
		if len(node.items) < btree.maxItems() {
			return zeroVal, nil, false
		}
		promotedItem, splitNode = btree.split(node)
		return promotedItem, splitNode, true
	}

	node.items.insertAt(promotedItem.key, promotedItem.value, idx)
	node.children.insertAt(splitNode, idx+1)
	if len(node.items) < btree.maxItems() {
		return zeroVal, nil, false
	}
	promotedItem, splitNode = btree.split(node)
	return promotedItem, splitNode, true

}

func (btree *BTree[K, V]) Delete(k K) bool {
	root := btree.root
	if root == nil {
		return false
	}
	return btree.delete(k, btree.root)
}

func (btree *BTree[K, V]) delete(k K, node *Node[K, V]) bool {
	idx, found := node.items.find(k)
	if found {
		if node.isLeaf() {
			node.items.deleteAt(idx)
			return true
		}

		if predecessors := node.children[idx]; len(predecessors.items) > btree.minItems() {
			promotedItem := predecessors.items[len(predecessors.items)-1]
			btree.delete(promotedItem.key, predecessors)
			node.items.deleteAt(idx)
			node.items.insertAt(promotedItem.key, promotedItem.value, idx)
		} else if succesors := node.children[idx+1]; len(succesors.items) > btree.minItems() {
			promotedItem := succesors.items[0]
			btree.delete(promotedItem.key, succesors)
			node.items.deleteAt(idx)
			node.items.insertAt(promotedItem.key, promotedItem.value, idx)
		} else {
			// TODO: Check if node has mutiple keys
			predecessors.items = append(predecessors.items, node.items[idx])
			predecessors.items = append(predecessors.items, succesors.items...)

			if !predecessors.isLeaf() {
				predecessors.children = append(predecessors.children, succesors.children...)
			}
			deletedItem := node.items.deleteAt(idx)
			node.children.deleteAt(idx + 1)
			btree.delete(deletedItem.key, predecessors)
		}

		return true
	}

	if node.isLeaf() {
		return false
	}

	child := node.children[idx]
	if len(child.items) > btree.minItems() {
		return btree.delete(k, child)
	}

	hasLeftSibling := idx > 0
	hasRightSibling := idx < len(node.children)-1

	if hasLeftSibling && len(node.children[idx-1].items) > btree.minItems() {
		sibling := node.children[idx-1]
		demotedItem := node.items[idx-1]
		child.items.insertAt(demotedItem.key, demotedItem.value, 0)
		if !sibling.isLeaf() {
			siblingChild := sibling.children.deleteAt(len(sibling.children) - 1)
			child.children.insertAt(siblingChild, 0)
		}
		promotedItem := sibling.items.deleteAt(len(sibling.items) - 1)
		node.items[idx-1] = promotedItem
	} else if hasRightSibling && len(node.children[idx+1].items) > btree.minItems() {
		sibling := node.children[idx+1]
		child.items = append(child.items, node.items[idx])
		if !child.isLeaf() {
			child.children = append(child.children, sibling.children.deleteAt(0))
		}
		node.items[idx] = sibling.items.deleteAt(0)

	} else {
		if hasRightSibling {
			node.merge(idx)
		} else {
			node.merge(idx - 1)
		}

	}

	return btree.delete(k, child)

}

func (node *Node[K, V]) merge(i int) {
	child, sibling := node.children[i], node.children[i+1]

	child.items = append(child.items, node.items.deleteAt(i))
	child.items = append(child.items, sibling.items...)

	if !child.isLeaf() {
		child.children = append(child.children, sibling.children...)
	}
	node.children.deleteAt(i + 1)

}
