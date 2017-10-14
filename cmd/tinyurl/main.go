package main

import "github.com/mizoR/go-tinyurl"

func main() {
	s := tinyurl.NewServer()

	s.Start()
}
