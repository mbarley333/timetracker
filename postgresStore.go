package timetracker

import (
	"database/sql"
	"fmt"
)

const (
	BY_NAME             string = `SELECT task_name ,SUM(elapsed_time) elapsed_time FROM tasks WHERE task_name=$1 GROUP BY task_name`
	BY_SESSION          string = `SELECT id, task_name, start_time,elapsed_time FROM tasks t INNER JOIN task_session s ON t.id=s.taskid`
	INSERT              string = `INSERT INTO tasks(task_name, start_time) VALUES($1, $2) RETURNING id`
	REPORT              string = `SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name ORDER BY SUM(elapsed_time) DESC`
	LATEST              string = `SELECT task_name, start_time, elapsed_time FROM tasks ORDER BY start_time DESC LIMIT 10`
	UPDATE_STOPPED      string = `UPDATE tasks t SET elapsed_time=$1 FROM task_session s WHERE t.id = s.taskid`
	DELETE              string = `DELETE FROM tasks WHERE id=$1`
	UPSERT_TASK_SESSION string = `INSERT INTO task_session (username,taskid) VALUES ($1,$2) ON CONFLICT (username) DO UPDATE SET taskid=$2`
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

	stmt, err := p.Db.Prepare(INSERT)
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

	_, err := p.Db.Exec(UPSERT_TASK_SESSION, username, task.Id)
	if err != nil {
		return fmt.Errorf("unable to upsert task_session: %s", err)
	}
	return nil

}

func (p *PostgresStore) UpdateStopped(task Task) error {

	_, err := p.Db.Exec(UPDATE_STOPPED, task.ElapsedTimeSec)
	if err != nil {
		return fmt.Errorf("unable to update elapsed time: %s", err)
	}
	return nil

}

func (p *PostgresStore) Delete(task Task) error {

	_, err := p.Db.Exec(DELETE, task.Id)
	if err != nil {
		return fmt.Errorf("unable to delete record: %s", err)
	}
	return nil

}

// stop here
func (p *PostgresStore) GetTaskByName(taskname string) (Task, error) {

	rows, err := p.Db.Query(BY_NAME, taskname)
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

	rows, err := p.Db.Query(BY_SESSION)
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

	rows, err := p.Db.Query(REPORT)
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

	rows, err := p.Db.Query(LATEST)
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
