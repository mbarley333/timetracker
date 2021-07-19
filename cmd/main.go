package main

import (
	"log"
	"timetracker"
)

func main() {

	s := timetracker.NewServer()
	log.Fatal(s.ListenAndServe())

}
