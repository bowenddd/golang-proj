package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// define of hash function

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash function
	replicas int            // times of virtue nodes
	keys     []int          // a sorted slice to store keys
	hashMap  map[int]string // map of virtue node to true node, key is virtue node
}

func New(replicas int, hash Hash) *Map {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
