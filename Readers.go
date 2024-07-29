package main

import (
	"fmt"
	"io"
	"strings"
)

func concater() func(string) string {
	doc := ""
	return func(s string) string {
		doc +=  s + " "
		return doc
	}
}

func main() {
	r := strings.NewReader("l")
	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}