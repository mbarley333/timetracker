package main

import (
	"log"
	"timetracker"
)

func main() {

	// host, err := timetracker.GetEnvironmentVariable("TIMETRACKER_DB_HOST")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// port, err := timetracker.GetEnvironmentVariable("TIMETRACKER_DB_PORT")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// user, err := timetracker.GetEnvironmentVariable("TIMETRACKER_DB_USER")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// dbname, err := timetracker.GetEnvironmentVariable("TIMETRACKER_DB_NAME")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// convertPort, err := strconv.Atoi(port)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
	// 	host, convertPort, user, dbname)
	// fmt.Println(psqlInfo)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// env := timetracker.Env{Db: db}

	// test := timetracker.Task{
	// 	Name:        "hmm",
	// 	StartTime:   time.Now(),
	// 	ElapsedTime: 10.0,
	// }

	// err = env.Create(test)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// db, err := timetracker.ConnectDB()
	// if err != nil {
	// 	log.Fatalf("error connecting to database: %s", err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// env := &timetracker.Env{Db: db}

	// test := timetracker.Task{
	// 	Name:        "hmm",
	// 	StartTime:   time.Now(),
	// 	ElapsedTime: 10.0,
	// }

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// h

	s := timetracker.NewServer()
	log.Fatal(s.ListenAndServe())

}
