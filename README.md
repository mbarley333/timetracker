# timetracker
The timetracker application allows you to track time spent on tasks.  The timetracker frontend and backend is built with Go and the data store is either SQLite or Postgres.  

## prerequisites
* Go 1.16
* Docker
* docker-compose 

## startup options



**1) containerized timetracker application and Postgres container**
<br>clone repo locally
```bash
docker-compose up
browse to: http://127.0.0.1:4000/home
```

**2) timetracker with SQLite**
<br>clone repo locally
```bash
cd to directory
go run ./cmd/main.go
browse to: http://127.0.0.1:4000/home
```
-----

-----

**3) timetracker with Postgres container**
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
go run ../../cmd/main.go
browse to: http://127.0.0.1:4000/home
```




## Goals
To learn and become more familiar with the following aspects of the Go language:
* testing
* HTTP Server
* templates
* interfaces
* embed
* Postgres
* SQLite


