package main

import (
	"log"
)

func main() {
	s := NewHTTPServer(":8080")
	log.Fatal(s.ListenAndServe())
}
