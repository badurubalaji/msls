// Package database provides database utilities including connection pooling,
// transaction management, and context propagation.
package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionFunc is a function that runs within a database transaction.
type TransactionFunc func(tx *gorm.DB) error

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// If the function panics, the transaction is rolled back and the panic is re-raised.
// The transaction is automatically committed if the function returns nil.
func WithTransaction(ctx context.Context, db *gorm.DB, fn TransactionFunc) error {
	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw the panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("rollback failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithNestedTransaction executes a function within a nested transaction (savepoint).
// If the function returns an error, the savepoint is rolled back.
// This is useful for partial rollbacks within a larger transaction.
func WithNestedTransaction(tx *gorm.DB, fn TransactionFunc) error {
	// Begin a savepoint
	savepoint := tx.SavePoint("sp1")
	if savepoint.Error != nil {
		return fmt.Errorf("failed to create savepoint: %w", savepoint.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.RollbackTo("sp1")
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.RollbackTo("sp1").Error; rbErr != nil {
			return fmt.Errorf("rollback to savepoint failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	return nil
}

// ExecuteInTransaction is a generic helper that executes a function and returns a result.
// This is useful when you need to return a value from within a transaction.
func ExecuteInTransaction[T any](ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) (T, error)) (T, error) {
	var result T

	err := WithTransaction(ctx, db, func(tx *gorm.DB) error {
		var fnErr error
		result, fnErr = fn(tx)
		return fnErr
	})

	return result, err
}

// TxFromContext retrieves a transaction from the context, or returns the original db if not found.
func TxFromContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db.WithContext(ctx)
}

// ContextWithTx adds a transaction to the context.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	// txKey is the context key for database transactions.
	txKey contextKey = "db_tx"
)
