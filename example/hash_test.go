package test

import (
	"github.com/R-Goys/RigoCache/pkg/consistenthash"
	"strconv"
	"testing"
)

func Test_Hash(t *testing.T) {
	//测试了好几组数据，所有数据都是绑在一个节点上的，才发现是这个哈希算法太拉了
	m := consistenthash.New(1000, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	m.Add("6", "4", "2")
	testCases := map[string]string{
		"1516512": "测试1",
		"3":       "测试2",
		"242":     "测试3",
		"6":       "测试4",
	}
	for k := range testCases {
		t.Logf("Asking for %s, got %s", k, m.Get(k))
	}
	m.Add("15")

	for k := range testCases {
		t.Logf("Asking for %s, got %s", k, m.Get(k))
	}
}
