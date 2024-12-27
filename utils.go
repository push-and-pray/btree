package main

import "sort"

func (btree *BTree[K, V]) insertOrdered(slice []Element[K, V], k K, v V) ([]Element[K, V], int) {
	// Find the index of position where key fits
	index := sort.Search(len(slice), func(i int) bool {
		return slice[i].key > k
	})

	// Grow slice by one, shift [index:] to right(making space for new element)
	// and insert it where it belongs
	slice = append(slice, Element[K, V]{})
	copy(slice[index+1:], slice[index:])
	slice[index] = Element[K, V]{k, v}
	return slice, index
}
