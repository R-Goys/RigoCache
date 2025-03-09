package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 哈希函数
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //哈希函数
	replicas int            //虚拟节点倍数，也就是每一个真实节点有多少个虚拟节点
	keys     []int          //哈希环
	hashMap  map[int]string //虚拟节点和真实节点映射
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
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

func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		delete(m.hashMap, hash)
		for index, val := range m.keys {
			if val == hash {
				m.keys = append(m.keys[:index], m.keys[index+1:]...)
				break
			}
		}
	}

	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//计算哈希值
	hash := int(m.hash([]byte(key)))
	//找到第一个大于hash值的索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//返回应该从哪个节点取数据
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
