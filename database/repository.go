package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrWalletNotFound = errors.New("Wallet not found")
var ErrInsufficientFunds = errors.New("Isnsufficient funds")

type WalletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{
		pool: pool,
	}
}

func (r *WalletRepository) Create(ctx context.Context) (*Wallet, error) {
	sqlQuery := `
	INSERT INTO wallet(balance) 
	VALUES (0) 
	RETURNING wallet_id, balance;
	`

	var w Wallet
	err := r.pool.QueryRow(ctx, sqlQuery).Scan(&w.WalletID, &w.Balance)
	if err != nil {
		return nil, fmt.Errorf("Create wallet: %w", err)
	}

	return &w, nil
}

func (r *WalletRepository) GetBalance(ctx context.Context, waletID uuid.UUID) (*Wallet, error) {
	sqlQuery := `
	SELECT wallet_id, balance 
	FROM wallet 
	WHERE wallet_id = $1;
	`

	var wallet Wallet
	err := r.pool.QueryRow(ctx, sqlQuery, waletID).Scan(&wallet.WalletID, &wallet.Balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWalletNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("Get balance: %w", err)
	}

	return &wallet, nil
}

func (r *WalletRepository) Deposit(ctx context.Context, walletID uuid.UUID, amount float64) (*Wallet, error) {
	return r.changeBalance(ctx, walletID, amount)
}

func (r *WalletRepository) Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) (*Wallet, error) {
	return r.changeBalance(ctx, walletID, -amount)
}

func (r *WalletRepository) changeBalance(ctx context.Context, walletID uuid.UUID, amount float64) (*Wallet, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	sqlSelectQuery := `
	SELECT wallet_id, balance
	FROM wallet
	WHERE wallet_id = $1 FOR UPDATE;
	`

	var wallet Wallet
	err = tx.QueryRow(ctx, sqlSelectQuery, walletID).Scan(&wallet.WalletID, &wallet.Balance)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWalletNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("Select for update: %w", err)
	}

	if wallet.Balance+amount < 0 {
		return nil, ErrInsufficientFunds
	}

	sqlUpdateQuery := `
	UPDATE wallet
	SET balance = balance + $1
	WHERE wallet_id = $2
	RETURNING wallet_id, balance
	`

	var updWallet Wallet
	err = tx.QueryRow(ctx, sqlUpdateQuery, amount, walletID).Scan(&updWallet.WalletID, &updWallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("Update balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("Commit tx: %w", err)
	}

	return &updWallet, nil
}
