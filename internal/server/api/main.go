package main

import (
	"fmt"
	"github.com/R-Goys/RigoCache/internal/server/api/Initialize"
	"github.com/R-Goys/RigoCache/internal/server/api/handle"
	"net/http"
)

func main() {
	Initialize.InitPool()
	Initialize.InitETCD()
	fmt.Println("RigoCache is up and running")
	http.ListenAndServe(":8080", handle.Getters)
}
