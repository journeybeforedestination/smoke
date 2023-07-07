package main

import (
	"log"

	"github.com/journeybeforedestination/smoke/tester"
)

func main() {
	s := tester.NewTestServer(":8081")
	log.Fatal(s.ListenAndServe())
}
