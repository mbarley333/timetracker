CREATE DATABASE timetracker;



CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    task_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    elapsed_time NUMERIC NOT NULL
);