CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    elapsed_time NUMERIC DEFAULT 0
);


CREATE TABLE IF NOT EXISTS task_session(
    taskid INTEGER
);