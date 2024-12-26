package main

import (
	"fmt"
	"testing"
)

func newPopulatedBTree() *BTree[int, string] {

	layer_3_1 := &Node[int, string]{
		elements: []Element[int, string]{{59, "59"}, {61, "61"}},
	}
	layer_3_2 := &Node[int, string]{
		elements: []Element[int, string]{{71, "71"}, {79, "79"}, {83, "83"}},
	}
	layer_3_3 := &Node[int, string]{
		elements: []Element[int, string]{{97, "97"}, {101, "101"}},
	}

	layer_2_1 := &Node[int, string]{
		elements: []Element[int, string]{{5, "5"}, {7, "7"}},
	}
	layer_2_2 := &Node[int, string]{
		elements: []Element[int, string]{{17, "17"}, {23, "23"}},
	}
	layer_2_3 := &Node[int, string]{
		elements: []Element[int, string]{{31, "31"}, {37, "37"}},
	}

	layer_1_0 := &Node[int, string]{
		elements: []Element[int, string]{{11, "11"}, {29, "29"}},
		children: []*Node[int, string]{layer_2_1, layer_2_2, layer_2_3},
	}
	layer_1_1 := &Node[int, string]{
		elements: []Element[int, string]{{67, "67"}, {89, "89"}},
		children: []*Node[int, string]{layer_3_1, layer_3_2, layer_3_3},
	}

	root := &Node[int, string]{
		elements: []Element[int, string]{{43, "43"}},
		children: []*Node[int, string]{layer_1_0, layer_1_1},
	}

	return &BTree[int, string]{max: 4, root: root}

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
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Get(%d)", test.key), func(t *testing.T) {
			result, success := btree.Get(test.key)
			if result != test.expected_val || success != test.expected_success {
				t.Errorf("expected %s and %t, got %s and %t", test.expected_val, test.expected_success, result, success)
			}
		})
	}
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

	keys := []int{696969, 123456, 999999, 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			btree.Get(key)
		}
	}
}

func TestDebug(t *testing.T) {
	btree := newPopulatedBTree()
	btree.Add(42, "42")
	println(btree.String())
}
