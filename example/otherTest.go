package main

import (
	"fmt"
)

func main() {
	mapper := make(map[string][]byte, 0)
	mapper["gejin"] = []byte("123")
	fmt.Println(len(mapper))
}
