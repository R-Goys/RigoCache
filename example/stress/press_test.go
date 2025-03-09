package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_press(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			s := "http://localhost:8080/RigoCache/score/Sam"
			get, err := http.Get(s)
			if err != nil {
				return
			}
			defer get.Body.Close()
			getBytes, err := ioutil.ReadAll(get.Body)
			if err != nil {
				return
			}
			getString := string(getBytes)
			fmt.Println(getString)
		}()
	}
	select {}
}
