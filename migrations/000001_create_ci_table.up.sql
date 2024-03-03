BEGIN;

CREATE TYPE task_status_t AS ENUM ('SUCCESS', 'FAILED');
COMMENT ON TYPE task_status_t IS 'Статусы запроса';

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    name  VARCHAR(255),
    command VARCHAR(255),
    status task_status_t NOT NULL,
    exit_code INT NOT NULL,
    output TEXT,
    exec_time FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT localtimestamp
);
COMMENT ON TABLE tasks IS 'task list';
COMMENT ON COLUMN tasks.id IS 'ID of record';
COMMENT ON COLUMN tasks.name IS 'name of task';
COMMENT ON COLUMN tasks.command IS 'execution command';
COMMENT ON COLUMN tasks.status IS 'status';
COMMENT ON COLUMN tasks.exit_code IS 'exit code of command';
COMMENT ON COLUMN tasks.output IS 'output of command';
COMMENT ON COLUMN tasks.exec_time IS 'exec time of command';
COMMENT ON COLUMN tasks.created_at IS 'created at';


CREATE TABLE IF NOT EXISTS pipelines (
    id SERIAL PRIMARY KEY,
    repository VARCHAR(255) NOT NULL,
    pipeline_id VARCHAR(50) NOT NULL,
    task_id BIGINT REFERENCES tasks(id)
);
COMMENT ON TABLE pipelines IS 'pipeline list';
COMMENT ON COLUMN pipelines.id IS 'ID of record';
COMMENT ON COLUMN pipelines.repository IS 'repository';
COMMENT ON COLUMN pipelines.pipeline_id IS 'id of pipeline';
COMMENT ON COLUMN pipelines.task_id IS 'id of task';

COMMIT;