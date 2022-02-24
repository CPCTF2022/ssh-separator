package badger

import (
	"context"
	"errors"
	"fmt"

	ctxManager "github.com/CPCTF2022/ssh-separator/pkg/context"
	"github.com/dgraph-io/badger/v3"
)

type Transaction struct {
	db *DB
}

func NewTransaction(db *DB) *Transaction {
	return &Transaction{
		db: db,
	}
}

type transactionType int

const (
	read transactionType = iota
	write
	none
)

func (t *Transaction) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := t.db.DB.Update(func(txn *badger.Txn) error {
		ctx := context.WithValue(ctx, ctxManager.TransactionKey, txn)

		return fn(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

func (t *Transaction) RTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := t.db.DB.View(func(txn *badger.Txn) error {
		ctx := context.WithValue(ctx, ctxManager.TransactionKey, txn)

		return fn(ctx)
	})
	if err != nil {
		return fmt.Errorf("failed in transaction: %w", err)
	}

	return nil
}

func getTransaction(ctx context.Context) (*badger.Txn, error) {
	iTxn := ctx.Value(ctxManager.TransactionKey)
	if iTxn == nil {
		return nil, nil
	}

	txn, ok := iTxn.(*badger.Txn)
	if !ok {
		return nil, errors.New("broken transaction")
	}

	return txn, nil
}
