package dao

import (
	"context"
	"github.com/shopspring/decimal"
	"testing"
)

func TestWalletManager_Transfer(t *testing.T) {
	type args struct {
		ctx          context.Context
		fromWalletID int64
		toWalletID   int64
		amount       decimal.Decimal
	}
	tests := []struct {
		name            string
		args            args
		setupBalances   map[int64]decimal.Decimal
		wantErr         bool
		wantFromBalance decimal.Decimal
		wantToBalance   decimal.Decimal
	}{
		{
			name: "正常转账",
			args: args{
				ctx:          context.Background(),
				fromWalletID: 1,
				toWalletID:   2,
				amount:       decimal.NewFromInt(100),
			},
			setupBalances: map[int64]decimal.Decimal{
				1: decimal.NewFromInt(150),
				2: decimal.NewFromInt(10),
			},
			wantErr:         false,
			wantFromBalance: decimal.NewFromInt(50),
			wantToBalance:   decimal.NewFromInt(110),
		},
		{
			name: "余额不足",
			args: args{
				ctx:          context.Background(),
				fromWalletID: 1,
				toWalletID:   2,
				amount:       decimal.NewFromInt(100),
			},
			setupBalances: map[int64]decimal.Decimal{
				1: decimal.NewFromInt(99),
				2: decimal.NewFromInt(0),
			},
			wantErr: true,
		},
		{
			name: "金额非法(<=0)",
			args: args{
				ctx:          context.Background(),
				fromWalletID: 1,
				toWalletID:   2,
				amount:       decimal.Zero,
			},
			setupBalances: map[int64]decimal.Decimal{
				1: decimal.NewFromInt(100),
				2: decimal.NewFromInt(0),
			},
			wantErr: true,
		},
		{
			name: "自转账",
			args: args{
				ctx:          context.Background(),
				fromWalletID: 1,
				toWalletID:   1,
				amount:       decimal.NewFromInt(1),
			},
			setupBalances: map[int64]decimal.Decimal{
				1: decimal.NewFromInt(100),
			},
			wantErr: true,
		},
		{
			name: "收款钱包不存在",
			args: args{
				ctx:          context.Background(),
				fromWalletID: 1,
				toWalletID:   2,
				amount:       decimal.NewFromInt(1),
			},
			setupBalances: map[int64]decimal.Decimal{
				1: decimal.NewFromInt(100),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWalletManager()
			for id, bal := range tt.setupBalances {
				if _, err := w.Init(context.Background(), id); err != nil {
					t.Fatalf("Init(%d) err: %v", id, err)
				}
				w.mu.Lock()
				w.wallets[id].Balance = bal
				w.mu.Unlock()
			}

			err := w.Transfer(tt.args.ctx, tt.args.fromWalletID, tt.args.toWalletID, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			fromWallet, err := w.Get(context.Background(), tt.args.fromWalletID)
			if err != nil {
				t.Fatalf("Get(from) err: %v", err)
			}
			toWallet, err := w.Get(context.Background(), tt.args.toWalletID)
			if err != nil {
				t.Fatalf("Get(to) err: %v", err)
			}
			if !fromWallet.Balance.Equal(tt.wantFromBalance) {
				t.Fatalf("from balance=%s want=%s", fromWallet.Balance.String(), tt.wantFromBalance.String())
			}
			if !toWallet.Balance.Equal(tt.wantToBalance) {
				t.Fatalf("to balance=%s want=%s", toWallet.Balance.String(), tt.wantToBalance.String())
			}
		})
	}
}
