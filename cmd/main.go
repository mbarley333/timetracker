package main

import (
	"log"
	"timetracker"
)

func main() {

	s := timetracker.NewServer(
		timetracker.WithSqliteStore(),
	)
	log.Fatal(s.ListenAndServe())

}
