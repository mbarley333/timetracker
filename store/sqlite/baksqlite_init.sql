CREATE TABLE tasks (
    id INTEGER PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    elapsed_time NUMERIC DEFAULT 0
);


CREATE TABLE task_session(
    taskid INTEGER
);