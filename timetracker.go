package timetracker

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Env struct {
	Db *sql.DB
}

type Task struct {
	Name        string `db:"task_name"`
	Active      bool
	StartTime   time.Time `db:"start_time"`
	ElapsedTime float64   `db:"elasped_time"`
}

type Report struct {
	Task      string
	TotalTime float64
}

func NewTask(task string) Task {
	t := Task{
		Name: task,
	}
	return t
}

func (t Task) GetActive() bool {
	return t.Active
}

func (t Task) GetStartTime() time.Time {
	return t.StartTime
}

func (t *Task) Start(now time.Time) {
	t.Active = true
	t.StartTime = now
}

func (t *Task) Stop(now time.Time) {
	stop := now
	t.ElapsedTime = stop.Sub(t.StartTime).Seconds()
	t.Active = false

}

func (t Task) GetElapsedTime() float64 {
	return t.ElapsedTime
}

func (t Task) GetMessage() string {

	return fmt.Sprintf("You spent %.1f seconds on the %s task", t.ElapsedTime, t.Name)
}

// server
func ConnectDB() (*sql.DB, error) {

	host, err := GetEnvironmentVariable("TIMETRACKER_DB_HOST")
	if err != nil {
		log.Fatal(err)
	}
	port, err := GetEnvironmentVariable("TIMETRACKER_DB_PORT")
	if err != nil {
		log.Fatal(err)
	}
	user, err := GetEnvironmentVariable("TIMETRACKER_DB_USER")
	if err != nil {
		log.Fatal(err)
	}
	dbname, err := GetEnvironmentVariable("TIMETRACKER_DB_NAME")
	if err != nil {
		log.Fatal(err)
	}

	convertPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, convertPort, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetEnvironmentVariable(env string) (string, error) {

	value := os.Getenv(env)

	if value == "" {
		return "", fmt.Errorf("%s value not set", env)
	}
	return value, nil
}

// DB related
func (e *Env) Create(task Task) error {
	query, err := GenerateSQLQuery("insert")
	if err != nil {
		return fmt.Errorf("unable to generate insert SQL")
	}
	_, err = e.Db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating task in database: %s", err)
	}
	return nil
}

func GenerateSQLQuery(sql string) (string, error) {
	switch sql {
	case "insert":
		return `INSERT INTO tasks(task_name, start_time, elased_time) VALUES($1, $2, $3)`, nil
	case "report":
		return `SELECT task, SUM(elapsed_time) total_time FROM tasks GROUP BY task`, nil
	}
	return "", fmt.Errorf("unable to generate sql based on input paramter: %s", sql)
}

// func (e *Env) GetReport() ([]Report, error) {

// 	query, err := GenerateSQLQuery("report")
// 	if err != nil {
// 		return nil, fmt.Errorf("error: %s", err)
// 	}
// 	rows, err := e.Db.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get report: %s", err)
// 	}
// 	defer rows.Close()

// 	var reports []Report

// 	for rows.Next() {
// 		r := Report{}
// 		err := rows.Scan(r.Task, r.TotalTime)
// 		if err != nil {
// 			return []Report{}, fmt.Errorf("failed to scan Report: %s", err)
// 		}
// 		reports = append(reports, r)
// 	}

// 	return reports, nil

// }

func ParseRows(r *sql.Rows) ([]Report, error) {

	var reports []Report
	for r.Next() {
		var report Report
		if err := r.Scan(&report.TotalTime, &report.Task); err != nil {
			return []Report{}, fmt.Errorf("unable to scan report: %s", err)
		}
		reports = append(reports, report)
	}

	return reports, nil

}
