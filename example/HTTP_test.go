package test

import (
	Rigo "github.com/R-Goys/RigoCache/core"
	RigoHTTP "github.com/R-Goys/RigoCache/http"
	"log"
	"net/http"
	"testing"
)

func Test_http(t *testing.T) {
	var getter Rigo.GetterFunc = func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			log.Println("[slowDB] got", key)
			return []byte(v), nil
		}
		log.Println("[slowDB] key not found")
		return nil, nil
	}
	Rigo.NewGroup("score", getter, 7)
	addr := "localhost:8080"
	srv := RigoHTTP.NewHttpPool(addr)
	log.Println("[HTTP srv] start in", addr)
	log.Fatal(http.ListenAndServe(addr, srv))
}
