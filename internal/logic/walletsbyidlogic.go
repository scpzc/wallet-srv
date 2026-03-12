package logic

import (
	"context"

	"wallet-srv/internal/svc"
	"wallet-srv/wallet_srv"

	"github.com/zeromicro/go-zero/core/logx"
)

type WalletsByIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWalletsByIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletsByIDLogic {
	return &WalletsByIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WalletsByIDLogic) WalletsByID(req *wallet_srv.WalletsByIDReq) (*wallet_srv.WalletsByIDResp, error) {
	wallet, err := l.svcCtx.WalletManager.Get(l.ctx, req.WalletID)
	if err != nil {
		return nil, err
	}
	return &wallet_srv.WalletsByIDResp{
		WalletID: wallet.UserId,
		Balance:  wallet.Balance.String(),
	}, nil
}
