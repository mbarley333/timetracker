# timetracker
The timetracker application allows you to track time spent on tasks.  The project allows you to use either a SQLite or Postgres as the data store.

## prerequisites
* Go 1.16
* Docker
* docker-compose 

# Startup options

## containerized timetracker application and Postgres container
```bash
docker-compose up
```
Browse to: http://127.0.0.1/home

## timetracker with SQLite
```bash
go run cmd/main.go
```
Browse to: http://127.0.0.1/home

## timetracker with Postgres container
* edit cmd/main.go:

```bash
package main

import (
	"log"
	"timetracker"
)

func main() {

	conn := "host=localhost port=5432 user=postgres dbname=timetracker sslmode=disable"
	s := timetracker.NewServer(
		timetracker.WithPostgresStore(conn),
	)
	log.Fatal(s.ListenAndServe())

}
```


```bash
cd store/pg
docker-compose up
go run cmd/main.go
```
Browse to: http://127.0.0.1/home



## Goals
To learn and become more familiar with the following aspects of the Go language:
* testing
* HTTP Server
* templates
* interfaces
* Postgres
* SQLite


