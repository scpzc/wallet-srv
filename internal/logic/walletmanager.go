package logic

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
	wallets sync.Map
}

func NewWalletManager() *WalletManager {
	return &WalletManager{}
}

func (w *WalletManager) Init(ctx context.Context, walletID int64) (int64, error) {
	wallet := &Wallet{UserId: walletID}
	w.wallets.Store(walletID, wallet)
	return walletID, nil
}

func (w *WalletManager) Get(ctx context.Context, walletID int64) (*Wallet, error) {
	wallet, ok := w.wallets.Load(walletID)
	if !ok {
		logx.Errorw("wallet not found", logx.Field("wallet_id", walletID))
		return nil, nil
	}
	return wallet.(*Wallet), nil
}
func (w *WalletManager) Transfer(ctx context.Context, fromWalletID int64, toWalletID int64, amount decimal.Decimal) error {
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
