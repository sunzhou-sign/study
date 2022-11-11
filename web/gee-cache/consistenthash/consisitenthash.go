package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// 一致性哈希算法的主数据结构
type Map struct {
	Hash
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点与真实节点的映射表 键是虚拟节点的哈希值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		Hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.Hash == nil {
		m.Hash = crc32.ChecksumIEEE
	}

	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对每一个真实节点 key，对应创建 m.replicas 个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点
			hash := int(m.Hash([]byte(strconv.Itoa(i) + key)))
			// 添加到环上
			m.keys = append(m.keys, hash)
			// 在 hashMap 中增加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	// 环上的哈希值排序。
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.Hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
