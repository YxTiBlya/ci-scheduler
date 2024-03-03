package db

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/YxTiBlya/ci-scheduler/internal/service/dto"
)

func (db *DB) transactionEnd(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		if err = tx.Rollback(ctx); err != nil {
			db.log.Error().Err(err).Msg("failed transaction rollback")
		}
		return
	}

	if err = tx.Commit(ctx); err != nil {
		db.log.Error().Err(err).Msg("failed transaction commit")
	}
}

const qInsertTask = `
	INSERT INTO tasks (name, command, status, exit_code, output, exec_time)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
`
const qInsertPipeline = `
	INSERT INTO pipelines (repository, pipeline_id, task_id)
	VALUES ($1, $2, $3)
`

func (db *DB) InsertPipline(ctx context.Context, data *dto.TaskData) error {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		db.log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}
	defer db.transactionEnd(ctx, tx, err)

	var id int
	err = db.Pool.QueryRow(
		ctx, qInsertTask,
		data.Name,
		data.Command,
		data.Status,
		data.ExitCode,
		data.Output,
		data.ExecTime,
	).Scan(&id)
	if err != nil {
		db.log.Error().Err(err).Msg("failed to insert task")
		return err
	}

	if err = db.Pool.QueryRow(ctx, qInsertPipeline, data.Repo, data.PipelineID, id).Scan(); err != pgx.ErrNoRows {
		db.log.Error().Err(err).Msg("failed to insert pipeline")
		return err
	}

	return nil
}
