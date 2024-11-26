package datasource

import (
	"contentgit/foundation"
	"context"
	"fmt"

	"runtime/debug"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func TransactionalWithContext(ctx context.Context, fn func(ctx context.Context) error) error {
	db := foundation.ContextProvider().GetDB(ctx)
	if db == nil {
		return errors.New("gorm middleware setup required.")
	}

	tx := db.Begin()
	ctx = foundation.ContextProvider().SetDB(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("TransactionManager Error", zap.Any("recover", r))
			zap.L().Error(fmt.Sprintf("Recovered. Error: %v \n %v", r, string(debug.Stack())))
			tx.Rollback()
		}
	}()

	if err := fn(ctx); err != nil {
		if err := tx.Rollback().Error; err != nil {
			zap.L().Error("TransactionManager Rollback Error", zap.Error(err))
			return err
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		zap.L().Error("TransactionManager Commit Error", zap.Error(err))
		return err
	}
	return nil
}
