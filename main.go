package main

import (
	"log"
)

func main() {
	s := NewFhirServer(":8080")
	log.Fatal(s.ListenAndServe())
}
