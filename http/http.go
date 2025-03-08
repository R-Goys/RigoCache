package RigoHTTP

import (
	"fmt"
	Rigo "github.com/R-Goys/RigoCache/core"
	"log"
	"net/http"
	"strings"
)

var defaultbasepath string = "/RigoCache/"

type HttpPool struct {
	basePath string
	self     string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		basePath: defaultbasepath,
		self:     self,
	}
}
func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// 这里也相当于是实现了接口，url格式固定为:"/{basePath}/{GroupName}/{Key}"
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log(r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := Rigo.GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
