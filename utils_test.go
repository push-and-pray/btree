package main

import (
	"reflect"
	"testing"
)

func TestInsertOrdered(t *testing.T) {
	btree := &BTree[int, string]{}

	tests := []struct {
		name     string
		slice    []Element[int, string]
		k        int
		v        string
		expected []Element[int, string]
	}{
		{
			name:     "insert into empty slice",
			slice:    []Element[int, string]{},
			k:        1,
			v:        "a",
			expected: []Element[int, string]{{1, "a"}},
		},
		{
			name:     "insert into beginning",
			slice:    []Element[int, string]{{2, "b"}, {3, "c"}},
			k:        1,
			v:        "a",
			expected: []Element[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		},
		{
			name:     "insert into middle",
			slice:    []Element[int, string]{{1, "a"}, {3, "c"}},
			k:        2,
			v:        "b",
			expected: []Element[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		},
		{
			name:     "insert into end",
			slice:    []Element[int, string]{{1, "a"}, {2, "b"}},
			k:        3,
			v:        "c",
			expected: []Element[int, string]{{1, "a"}, {2, "b"}, {3, "c"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := btree.insertOrdered(tt.slice, tt.k, tt.v)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func BenchmarkInsertOrdered(b *testing.B) {
	btree := &BTree[int, string]{}

	benchmarks := []struct {
		name  string
		slice []Element[int, string]
		k     int
		v     string
	}{
		{
			name:  "insert into empty slice",
			slice: []Element[int, string]{},
			k:     1,
			v:     "a",
		},
		{
			name:  "insert into beginning",
			slice: []Element[int, string]{{2, "b"}, {3, "c"}},
			k:     1,
			v:     "a",
		},
		{
			name:  "insert into middle",
			slice: []Element[int, string]{{1, "a"}, {3, "c"}},
			k:     2,
			v:     "b",
		},
		{
			name:  "insert into end",
			slice: []Element[int, string]{{1, "a"}, {2, "b"}},
			k:     3,
			v:     "c",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				btree.insertOrdered(bm.slice, bm.k, bm.v)
			}
		})
	}
}
