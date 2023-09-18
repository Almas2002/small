package transaction

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func Finish(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return err
		}
		return err
	} else {
		if commitErr := tx.Commit(ctx); commitErr != nil {
			return err
		}
		return nil
	}
}
