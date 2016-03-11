package main

import (
	. "fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fi, err := ioutil.ReadFile("./token")
	if err != nil {
		panic(err)
	}

	Print(string(fi))

	resp, _ := http.Get("https://api.github.com")
	Print(resp)
}
