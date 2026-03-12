package dao

import (
	"context"
	"github.com/shopspring/decimal"
	"sync"
	"testing"
)

func TestWalletManager_Transfer(t *testing.T) {
	type fields struct {
		wallets sync.Map
	}
	type args struct {
		ctx          context.Context
		fromWalletID int64
		toWalletID   int64
		amount       decimal.Decimal
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "正常",
			fields: fields{
				wallets: sync.Map{},
			},
			args: args{
				ctx:          context.Background(),
				fromWalletID: 1,
				toWalletID:   2,
				amount:       decimal.NewFromInt(100),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WalletManager{
				wallets: tt.fields.wallets,
			}
			walletID1, err := w.Init(context.Background(), 1)
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("wallet1:%v", walletID1)
			walletID2, err := w.Init(context.Background(), 2)
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("wallet2:%v", walletID2)
			if err := w.Transfer(tt.args.ctx, tt.args.fromWalletID, tt.args.toWalletID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
