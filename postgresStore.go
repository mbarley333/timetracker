package timetracker

import (
	"database/sql"
	"fmt"
)

type PostgresStore struct {
	Db *sql.DB
}

func NewPostgresStore(conn string) (*PostgresStore, error) {

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{Db: db}, nil
}

// DB related
func (p *PostgresStore) Create(task Task) (int, error) {

	query, err := GenerateSQLQuery("insert")
	if err != nil {
		return 0, fmt.Errorf("unable to generate insert SQL")
	}
	stmt, err := p.Db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("unable to prepare query")
	}
	defer stmt.Close()

	var taskid int
	err = stmt.QueryRow(task.Name, task.StartTime).Scan(&taskid)

	if err != nil {
		return 0, fmt.Errorf("error creating task in database: %s", err)
	}
	return taskid, nil
}

func (p *PostgresStore) NewTaskSession(task Task) error {
	username := "app"
	query, err := GenerateSQLQuery("upsertTaskSession")
	if err != nil {
		return fmt.Errorf("unable to generate update SQL")
	}

	_, err = p.Db.Exec(query, username, task.Id)
	if err != nil {
		return fmt.Errorf("unable to upsert task_session: %s", err)
	}
	return nil

}

func (p *PostgresStore) UpdateStopped(task Task) error {

	query, err := GenerateSQLQuery("updateStopped")
	if err != nil {
		return fmt.Errorf("unable to generate update SQL")
	}
	_, err = p.Db.Exec(query, task.ElapsedTimeSec)
	if err != nil {
		return fmt.Errorf("unable to update elapsed time: %s", err)
	}
	return nil

}

func (p *PostgresStore) Delete(task Task) error {

	query, err := GenerateSQLQuery("delete")
	if err != nil {
		return fmt.Errorf("unable to generate update SQL")
	}

	_, err = p.Db.Exec(query, task.Id)
	if err != nil {
		return fmt.Errorf("unable to delete record: %s", err)
	}
	return nil

}

// stop here
func (p *PostgresStore) GetTaskByName(taskname string) (Task, error) {

	query, err := GenerateSQLQuery("byname")
	if err != nil {
		return Task{}, fmt.Errorf("error: %s", err)
	}
	rows, err := p.Db.Query(query, taskname)
	if err != nil {
		return Task{}, fmt.Errorf("failed to get report: %s", err)
	}
	defer rows.Close()

	task, err := ParseRowsTaskByName(rows)
	if err != nil {
		return Task{}, fmt.Errorf("failed to parse rows: %s", err)
	}

	return task, nil
}

func (p *PostgresStore) GetTaskBySession() (Task, error) {

	query, err := GenerateSQLQuery("bysession")
	if err != nil {
		return Task{}, fmt.Errorf("error: %s", err)
	}

	rows, err := p.Db.Query(query)
	if err != nil {
		return Task{}, fmt.Errorf("failed to get report: %s", err)
	}
	defer rows.Close()

	task, err := ParseRowsTask(rows)
	if err != nil {
		return Task{}, fmt.Errorf("failed to parse rows: %s", err)
	}

	return task, nil
}

func (p *PostgresStore) GetReport() ([]Report, error) {

	query, err := GenerateSQLQuery("report")
	if err != nil {
		return []Report{}, fmt.Errorf("error: %s", err)
	}
	rows, err := p.Db.Query(query)
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

func (p *PostgresStore) GetLatest() ([]Task, error) {

	query, err := GenerateSQLQuery("latest")
	if err != nil {
		return []Task{}, fmt.Errorf("error: %s", err)
	}
	rows, err := p.Db.Query(query)
	if err != nil {
		return []Task{}, fmt.Errorf("failed to get report: %s", err)
	}
	defer rows.Close()

	tasks, err := ParseRowsTasks(rows)
	if err != nil {
		return []Task{}, fmt.Errorf("failed to parse rows: %s", err)
	}

	return tasks, nil

}

func GenerateSQLQuery(sql string) (string, error) {
	switch sql {
	case "byname":
		return `SELECT task_name ,SUM(elapsed_time) elapsed_time FROM tasks WHERE task_name=$1 GROUP BY task_name`, nil
	case "bysession":
		return `SELECT id, task_name, start_time,elapsed_time FROM tasks t INNER JOIN task_session s ON t.id=s.taskid`, nil
	case "insert":
		return `INSERT INTO tasks(task_name, start_time) VALUES($1, $2) RETURNING id`, nil
	case "report":
		return `SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name ORDER BY SUM(elapsed_time) DESC`, nil
	case "latest":
		return `SELECT task_name, start_time, elapsed_time FROM tasks ORDER BY start_time DESC LIMIT 10`, nil
	case "updateStopped":
		return `UPDATE tasks t SET elapsed_time=$1 FROM task_session s WHERE t.id = s.taskid`, nil
	case "delete":
		return `DELETE FROM tasks WHERE id=$1`, nil
	case "upsertTaskSession":
		return `INSERT INTO task_session (username,taskid) VALUES ($1,$2) ON CONFLICT (username) DO UPDATE SET taskid=$2`, nil
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
	var task Task

	for r.Next() {

		if err := r.Scan(&task.Name, &task.StartTime, &task.ElapsedTimeSec); err != nil {
			return []Task{}, fmt.Errorf("unable to scan tasks: %s", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil

}

func ParseRowsTask(r *sql.Rows) (Task, error) {

	var task Task

	for r.Next() {

		if err := r.Scan(&task.Id, &task.Name, &task.StartTime, &task.ElapsedTimeSec); err != nil {
			return Task{}, fmt.Errorf("unable to scan tasks: %s", err)
		}

	}

	return task, nil

}

func ParseRowsTaskByName(r *sql.Rows) (Task, error) {

	var task Task

	for r.Next() {

		if err := r.Scan(&task.Name, &task.ElapsedTimeSec); err != nil {
			return Task{}, fmt.Errorf("unable to scan tasks: %s", err)
		}

	}

	return task, nil

}
