package dao

import (
	"context"
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
)

type IWallet interface {
	// 初始化钱包
	Init(ctx context.Context, walletID int64) (int64, error)
	// 获取钱包
	Get(ctx context.Context, walletID int64) (*Wallet, error)
	// 转账
	Transfer(ctx context.Context, fromWalletID int64, toWalletID int64, amount decimal.Decimal) error
}

type Wallet struct {
	UserId  int64           `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
}

type WalletManager struct {
	mu      sync.RWMutex
	wallets map[int64]*Wallet
}

func NewWalletManager() *WalletManager {
	return &WalletManager{
		wallets: make(map[int64]*Wallet),
	}
}

// 初始化钱包
func (w *WalletManager) Init(ctx context.Context, walletID int64) (int64, error) {
	if walletID <= 0 {
		return 0, ErrInvalidWalletID
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.wallets[walletID]; ok {
		return walletID, nil
	}
	w.wallets[walletID] = &Wallet{
		UserId:  walletID,
		Balance: decimal.Zero,
	}
	return walletID, nil
}

// 获取钱包
func (w *WalletManager) Get(ctx context.Context, walletID int64) (*Wallet, error) {
	if walletID <= 0 {
		return nil, ErrInvalidWalletID
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	wallet, ok := w.wallets[walletID]
	if !ok || wallet == nil {
		logx.Errorw("wallet not found", logx.Field("wallet_id", walletID))
		return nil, ErrWalletNotFound
	}
	// 返回副本，避免外部绕过锁修改内部状态
	return &Wallet{
		UserId:  wallet.UserId,
		Balance: wallet.Balance,
	}, nil
}

// 转账
func (w *WalletManager) Transfer(ctx context.Context, fromWalletID int64, toWalletID int64, amount decimal.Decimal) error {
	if fromWalletID <= 0 || toWalletID <= 0 {
		return ErrInvalidWalletID
	}
	if fromWalletID == toWalletID {
		return ErrSameWalletTransfer
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	fromWallet, ok := w.wallets[fromWalletID]
	if !ok || fromWallet == nil {
		logx.Errorw("wallet not found", logx.Field("wallet_id", fromWalletID))
		return ErrWalletNotFound
	}

	toWallet, ok := w.wallets[toWalletID]
	if !ok || toWallet == nil {
		logx.Errorw("wallet not found", logx.Field("wallet_id", toWalletID))
		return ErrWalletNotFound
	}

	// TODO 没有充值接口，所以暂时去掉这个判断
	//if fromWallet.Balance.LessThan(amount) {
	//	return ErrInsufficientBalance
	//}

	fromWallet.Balance = fromWallet.Balance.Sub(amount)
	toWallet.Balance = toWallet.Balance.Add(amount)

	logx.Infow("transfer success",
		logx.Field("from", fromWalletID),
		logx.Field("to", toWalletID),
		logx.Field("amount", amount.String()),
		logx.Field("from_balance", fromWallet.Balance.String()),
		logx.Field("to_balance", toWallet.Balance.String()),
	)
	return nil
}

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrInvalidWalletID     = errors.New("invalid wallet id")
	ErrSameWalletTransfer  = errors.New("same wallet transfer")
)
