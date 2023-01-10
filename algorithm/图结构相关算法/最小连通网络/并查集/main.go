package main

import "fmt"

// 简单版并查集 是否形成环
func main() {
	spots := []byte{'A', 'B', 'C', 'D'}
	sets := make(map[byte]*bool, len(spots))

	for _, v := range spots {
		sets[v] = new(bool)
		fmt.Printf("%c:%p\n", v, sets[v])
	}

	if !isSameSet(sets, 'A', 'B') {
		union(sets, 'A', 'B')
	} else {
		fmt.Println("A,B已经在一个集合!")
	}

	if !isSameSet(sets, 'B', 'C') {
		union(sets, 'B', 'C')
	} else {
		fmt.Println("B,C已经在一个集合!")
	}

	if !isSameSet(sets, 'D', 'B') {
		union(sets, 'D', 'B')
	} else {
		fmt.Println("A,B已经在一个集合!")
	}

	if !isSameSet(sets, 'A', 'D') {
		union(sets, 'A', 'D')
	} else {
		fmt.Println("A,D已经在一个集合!")
	}

}

// 判断是否在一个集合，地址是否相同
func isSameSet(sets map[byte]*bool, from, to byte) bool {
	return sets[from] != nil && sets[to] != nil && sets[from] == sets[to]
}

// 集合合并
func union(sets map[byte]*bool, from, to byte) {
	for k, _ := range sets {
		if sets[k] == sets[to] {
			sets[k] = sets[from]
		}
	}
}
