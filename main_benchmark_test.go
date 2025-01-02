package btree

import (
	"fmt"
	"math/rand"
	"testing"
)

func generateRandomKVPairs(n int) []struct {
	key   int
	value int
} {
	kvPairs := make([]struct {
		key   int
		value int
	}, n)
	for i := range n {
		kvPairs[i] = struct {
			key   int
			value int
		}{key: rand.Int(), value: rand.Int()}
	}
	return kvPairs
}

func BenchmarkBTreeInsert(b *testing.B) {
	testSizes := []int{1000, 10000, 100000}
	degrees := []int{2, 4, 8}

	for _, degree := range degrees {
		for _, size := range testSizes {
			b.Run(
				fmt.Sprintf("Insert_Degree_%v_Size_+%v", degree, size),
				func(b *testing.B) {
					btree := NewBtree[int, int](degree)
					kvPairs := generateRandomKVPairs(size)

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						for _, pair := range kvPairs {
							btree.Insert(pair.key, pair.value)
						}
					}
				},
			)
		}
	}
}

func BenchmarkBTreeGet(b *testing.B) {
	testSizes := []int{1000, 10000, 100000}
	degrees := []int{2, 4, 8}

	for _, degree := range degrees {
		for _, size := range testSizes {
			b.Run(
				fmt.Sprintf("Get_Degree_%v_Size_+%v", degree, size),
				func(b *testing.B) {
					btree := NewBtree[int, int](degree)
					kvPairs := generateRandomKVPairs(size)

					for _, pair := range kvPairs {
						btree.Insert(pair.key, pair.value)
					}

					keys := make([]int, size)
					for i, pair := range kvPairs {
						keys[i] = pair.key
					}

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						for _, key := range keys {
							btree.Get(key)
						}
					}
				},
			)
		}
	}
}

func BenchmarkBTreeDelete(b *testing.B) {
	testSizes := []int{1000, 10000, 100000}
	degrees := []int{2, 4, 8}

	for _, degree := range degrees {
		for _, size := range testSizes {
			b.Run(
				fmt.Sprintf("Delete_Degree_%v_Size_+%v", degree, size),
				func(b *testing.B) {
					btree := NewBtree[int, int](degree)
					kvPairs := generateRandomKVPairs(size)

					for _, pair := range kvPairs {
						btree.Insert(pair.key, pair.value)
					}

					keys := make([]int, size)
					for i, pair := range kvPairs {
						keys[i] = pair.key
					}

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						for _, key := range keys {
							btree.Delete(key)
						}
					}
				},
			)
		}
	}
}
