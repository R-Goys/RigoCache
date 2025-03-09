package test

import (
	Rigo2 "github.com/R-Goys/RigoCache/internal/core"
	"github.com/R-Goys/RigoCache/internal/http"
	"log"
	"net/http"
	"testing"
)

func Test_http(t *testing.T) {
	var getter Rigo2.GetterFunc = func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			log.Println("[slowDB] got", key)
			return []byte(v), nil
		}
		log.Println("[slowDB] key not found")
		return nil, nil
	}
	Rigo2.NewGroup("score", getter, 7)
	addr := "localhost:8080"
	srv := RigoHTTP.NewHttpPool(addr)
	log.Println("[HTTP srv] start in", addr)
	log.Fatal(http.ListenAndServe(addr, srv))
}
