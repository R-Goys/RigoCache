package RigoHTTP

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HttpGetter 实现了拿取KV键值对的接口，在流程中作为手去拿取数据
type HttpGetter struct {
	baseURL string
}

// Get 发送请求，拿取数据
func (h *HttpGetter) Get(Group, key string) ([]byte, error) {
	u := fmt.Sprintf("%s/%s/%s", h.baseURL, Group, key)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP GET returned status code %d", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}
