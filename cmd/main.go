package main

import (
	"log"
	"timetracker"
)

func main() {
	db, err := timetracker.ConnectDB()
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	env := &timetracker.Env{Db: db}
}
