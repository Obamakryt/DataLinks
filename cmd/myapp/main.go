package main

import "fmt"

func main() {

	type s struct {
	}
	S(s)
}
func S(err error) {
	fmt.Println(err)
} n