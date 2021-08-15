package timetracker

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	SQLByName            string = `SELECT task_name ,SUM(elapsed_time) elapsed_time FROM tasks WHERE task_name=$1 GROUP BY task_name`
	SQLBySession         string = `SELECT id, task_name, start_time,elapsed_time FROM tasks t INNER JOIN task_session s ON t.id=s.taskid`
	SQLInsert            string = `INSERT INTO tasks(task_name, start_time) VALUES($1, $2) RETURNING id`
	SQLReport            string = `SELECT task_name, SUM(elapsed_time) total_time FROM tasks GROUP BY task_name ORDER BY SUM(elapsed_time) DESC`
	SQLLatestTasks       string = `SELECT task_name, start_time, elapsed_time FROM tasks ORDER BY start_time DESC LIMIT 10`
	SQLUpdateStopped     string = `UPDATE tasks SET elapsed_time=$1 FROM task_session  WHERE tasks.id = task_session.taskid`
	SQLDelete            string = `DELETE FROM tasks WHERE id=$1`
	SQLInsertTaskSession string = `INSERT INTO task_session (taskid) VALUES ($1)`
	SQLDeleteTaskSession string = `DELETE FROM task_session`
)

type DBStore struct {
	Db *sql.DB
}

func NewPostgresStore(conn string) (*DBStore, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return &DBStore{Db: db}, nil
}

func NewSqliteStore() (*DBStore, error) {
	db, err := sql.Open("sqlite3", "./timetracker.db")
	if err != nil {
		return nil, err
	}

	loadSQLFile(db, "store/sqlite/sqlite_init.sql")

	return &DBStore{Db: db}, nil
}

func loadSQLFile(db *sql.DB, sqlFile string) error {
	file, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()
	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DBStore) Create(task Task) (int, error) {

	stmt, err := d.Db.Prepare(SQLInsert)
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

func (d *DBStore) NewTaskSession(task Task) error {

	_, err := d.Db.Exec(SQLDeleteTaskSession)
	if err != nil {
		return fmt.Errorf("unable to delete record: %s", err)
	}

	_, err = d.Db.Exec(SQLInsertTaskSession, task.Id)
	if err != nil {
		return fmt.Errorf("unable to insert task_session: %s", err)
	}

	return nil

}

func (d *DBStore) UpdateStopped(task Task) error {

	_, err := d.Db.Exec(SQLUpdateStopped, task.ElapsedTimeSec)
	if err != nil {
		return fmt.Errorf("unable to update elapsed time: %s", err)
	}
	return nil

}

func (d *DBStore) Delete(task Task) error {

	_, err := d.Db.Exec(SQLDelete, task.Id)
	if err != nil {
		return fmt.Errorf("unable to delete record: %s", err)
	}
	return nil

}

// stop here
func (d *DBStore) GetTaskByName(taskname string) (Task, error) {

	rows, err := d.Db.Query(SQLByName, taskname)
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

func (d *DBStore) GetTaskBySession() (Task, error) {

	rows, err := d.Db.Query(SQLBySession)
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

func (d *DBStore) GetReport() ([]Report, error) {

	rows, err := d.Db.Query(SQLReport)
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

func (d *DBStore) GetLatest() ([]Task, error) {

	rows, err := d.Db.Query(SQLLatestTasks)
	if err != nil {
		return []Task{}, fmt.Errorf("failed to get latest: %s", err)
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
