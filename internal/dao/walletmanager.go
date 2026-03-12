package dao

import (
	"context"
	"errors"
	"fmt"
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
	sync.RWMutex
	wallets sync.Map
}

func NewWalletManager() *WalletManager {
	return &WalletManager{}
}

// 初始化钱包
func (w *WalletManager) Init(ctx context.Context, walletID int64) (int64, error) {
	w.Lock()
	defer w.Unlock()
	wallet := &Wallet{UserId: walletID}
	w.wallets.Store(walletID, wallet)
	wallet11, ok := w.wallets.Load(walletID)
	fmt.Println(wallet11, ok)
	return walletID, nil
}

// 获取钱包
func (w *WalletManager) Get(ctx context.Context, walletID int64) (*Wallet, error) {
	w.RLock()
	defer w.RUnlock()
	wallet, ok := w.wallets.Load(walletID)
	if !ok {
		logx.Errorw("wallet not found", logx.Field("wallet_id", walletID))
		return nil, errors.New("wallet not found")
	}
	return wallet.(*Wallet), nil
}

// 转账
func (w *WalletManager) Transfer(ctx context.Context, fromWalletID int64, toWalletID int64, amount decimal.Decimal) error {
	w.Lock()
	defer w.Unlock()
	fromWallet, ok := w.wallets.Load(fromWalletID)
	if !ok {
		logx.Errorw("wallet not found", logx.Field("wallet_id", fromWalletID))
		return errors.New("wallet not found")
	}
	toWallet, ok := w.wallets.Load(toWalletID)
	if !ok {
		logx.Errorw("wallet not found", logx.Field("wallet_id", fromWalletID))
		return errors.New("wallet not found")
	}
	fromBalance := fromWallet.(*Wallet).Balance

	// TODO 因为没有加余额的方法，所以这里先不考虑
	//if fromBalance.LessThan(amount) {
	//	return errors.New("insufficient balance")
	//}

	toBalance := toWallet.(*Wallet).Balance

	fromWallet.(*Wallet).Balance = fromBalance.Sub(amount)
	toWallet.(*Wallet).Balance = toBalance.Add(amount)

	logx.Infow("transfer success",
		logx.Field("from", fromWalletID),
		logx.Field("to", toWalletID),
		logx.Field("amount", amount.String()),
		logx.Field("from_balance", fromWallet.(*Wallet).Balance.String()),
		logx.Field("to_balance", toWallet.(*Wallet).Balance.String()),
	)
	return nil
}
