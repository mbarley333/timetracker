package timetracker

import (
	"database/sql"
	"fmt"
)

type Env struct {
	Db *sql.DB
}

// DB related
func (e *Env) Create(task Task) error {

	fmt.Println(task)
	query, err := GenerateSQLQuery("insert")
	if err != nil {
		return fmt.Errorf("unable to generate insert SQL")
	}
	fmt.Println(query)
	_, err = e.Db.Exec(`INSERT INTO tasks(task_name, start_time, elapsed_time) VALUES($1, $2, $3)`, task.Name, task.StartTime, task.ElapsedTime)
	if err != nil {
		return fmt.Errorf("error creating task in database: %s", err)
	}
	return nil
}

func (e *Env) GetReport() ([]Report, error) {

	query, err := GenerateSQLQuery("report")
	if err != nil {
		return []Report{}, fmt.Errorf("error: %s", err)
	}
	rows, err := e.Db.Query(query)
	if err != nil {
		return []Report{}, fmt.Errorf("failed to get report: %s", err)
	}
	defer rows.Close()

	reports, err := ParseRowsReport(rows)
	if err != nil {
		return []Report{}, fmt.Errorf("failed to parse rows: %s", err)
	}

	return reports, nil

}

func (e *Env) GetLatest() ([]Task, error) {

	query, err := GenerateSQLQuery("latest")
	if err != nil {
		return []Task{}, fmt.Errorf("error: %s", err)
	}
	rows, err := e.Db.Query(query)
	if err != nil {
		return []Task{}, fmt.Errorf("failed to get report: %s", err)
	}
	defer rows.Close()

	reports, err := ParseRowsTasks(rows)
	if err != nil {
		return []Task{}, fmt.Errorf("failed to parse rows: %s", err)
	}

	return reports, nil

}

func GenerateSQLQuery(sql string) (string, error) {
	switch sql {
	case "insert":
		return `INSERT INTO tasks(task_name, start_time, elapsed_time) VALUES($1, $2, $3)`, nil
	case "report":
		return `SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name`, nil
	case "latest":
		return `SELECT task_name, start_time, elapsed_time FROM tasks ORDER BY start_time DESC LIMIT 10`, nil
	}
	return "", fmt.Errorf("unable to generate sql based on input paramter: %s", sql)
}

func ParseRowsReport(r *sql.Rows) ([]Report, error) {

	var reports []Report
	for r.Next() {
		var report Report
		if err := r.Scan(&report.Task, &report.TotalTime); err != nil {
			return []Report{}, fmt.Errorf("unable to scan report: %s", err)
		}
		reports = append(reports, report)
	}

	return reports, nil

}

func ParseRowsTasks(r *sql.Rows) ([]Task, error) {

	var tasks []Task
	for r.Next() {
		var task Task
		if err := r.Scan(&task.Name, &task.StartTime, &task.ElapsedTime); err != nil {
			return []Task{}, fmt.Errorf("unable to scan tasks: %s", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil

}
