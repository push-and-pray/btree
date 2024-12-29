package btree

import (
	"testing"
)

func TestFind(t *testing.T) {
	items := items[int, string]{
		{key: 1, value: "one"},
		{key: 3, value: "three"},
		{key: 5, value: "five"},
	}

	tests := []struct {
		key      int
		expected int
		found    bool
	}{
		{key: 0, expected: 0, found: false},
		{key: 1, expected: 0, found: true},
		{key: 2, expected: 1, found: false},
		{key: 3, expected: 1, found: true},
		{key: 4, expected: 2, found: false},
		{key: 5, expected: 2, found: true},
		{key: 6, expected: 3, found: false},
	}

	for _, test := range tests {
		idx, found := items.find(test.key)
		if idx != test.expected || found != test.found {
			t.Errorf("find(%d) = (%d, %v); expected (%d, %v)", test.key, idx, found, test.expected, test.found)
		}
	}
}

func TestInsertAt(t *testing.T) {
	testItems := items[int, string]{
		{key: 1, value: "one"},
		{key: 3, value: "three"},
		{key: 4, value: "four"},
	}

	tests := []struct {
		key      int
		value    string
		index    int
		expected items[int, string]
	}{
		{
			key:   2,
			value: "two",
			index: 1,
			expected: items[int, string]{
				{key: 1, value: "one"},
				{key: 2, value: "two"},
				{key: 3, value: "three"},
				{key: 4, value: "four"},
			},
		},
		{
			key:   5,
			value: "five",
			index: 4,
			expected: items[int, string]{
				{key: 1, value: "one"},
				{key: 2, value: "two"},
				{key: 3, value: "three"},
				{key: 4, value: "four"},
				{key: 5, value: "five"},
			},
		},
		{
			key:   0,
			value: "zero",
			index: 0,
			expected: items[int, string]{
				{key: 0, value: "zero"},
				{key: 1, value: "one"},
				{key: 2, value: "two"},
				{key: 3, value: "three"},
				{key: 4, value: "four"},
				{key: 5, value: "five"},
			},
		},
	}

	for _, test := range tests {
		testItems.insertAt(test.key, test.value, test.index)
		for i, item := range testItems {
			if item.key != test.expected[i].key || item.value != test.expected[i].value {
				t.Errorf("insertAt(%d, %s, %d) = %v; expected %v", test.key, test.value, test.index, testItems, test.expected)
				break
			}
		}
	}
}
func TestDeleteAt(t *testing.T) {
	testItems := items[int, string]{
		{key: 1, value: "one"},
		{key: 2, value: "two"},
		{key: 3, value: "three"},
		{key: 4, value: "four"},
	}

	tests := []struct {
		index    int
		expected items[int, string]
	}{
		{
			index: 1,
			expected: items[int, string]{
				{key: 1, value: "one"},
				{key: 3, value: "three"},
				{key: 4, value: "four"},
			},
		},
		{
			index: 0,
			expected: items[int, string]{
				{key: 3, value: "three"},
				{key: 4, value: "four"},
			},
		},
		{
			index: 1,
			expected: items[int, string]{
				{key: 3, value: "three"},
			},
		},
	}

	for _, test := range tests {
		testItems.deleteAt(test.index)
		if len(testItems) != len(test.expected) {
			t.Errorf("deleteAt(%d) = %v; expected %v", test.index, testItems, test.expected)
			continue
		}
		for i, item := range testItems {
			if item.key != test.expected[i].key || item.value != test.expected[i].value {
				t.Errorf("deleteAt(%d) = %v; expected %v", test.index, testItems, test.expected)
				break
			}
		}
	}
}

func TestDeleteAtChildren(t *testing.T) {
	child1 := &Node[int, string]{}
	child2 := &Node[int, string]{}
	child3 := &Node[int, string]{}
	child4 := &Node[int, string]{}

	testChildren := children[int, string]{child1, child2, child3, child4}

	tests := []struct {
		index    int
		expected children[int, string]
	}{
		{
			index:    1,
			expected: children[int, string]{child1, child3, child4},
		},
		{
			index:    0,
			expected: children[int, string]{child3, child4},
		},
		{
			index:    1,
			expected: children[int, string]{child3},
		},
	}

	for _, test := range tests {
		testChildren.deleteAt(test.index)
		if len(testChildren) != len(test.expected) {
			t.Errorf("deleteAt(%d) = %v; expected %v", test.index, testChildren, test.expected)
			continue
		}
		for i, child := range testChildren {
			if child != test.expected[i] {
				t.Errorf("deleteAt(%d) = %v; expected %v", test.index, testChildren, test.expected)
				break
			}
		}
	}
}
