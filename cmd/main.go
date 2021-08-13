package main

import (
	"log"
	"timetracker"
)

func main() {

	conn := "host=postgres port=5432 user=postgres dbname=timetracker sslmode=disable"
	s := timetracker.NewServer(
		timetracker.WithDataStore(conn),
	)
	log.Fatal(s.ListenAndServe())

}
