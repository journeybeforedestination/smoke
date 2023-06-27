package main

import (
	"log"

	"github.com/journeybeforedestination/smoke/launcher"
)

func main() {
	s := launcher.NewFhirServer(":8080")
	log.Fatal(s.ListenAndServe())
}
