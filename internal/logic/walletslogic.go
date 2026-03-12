package logic

import (
	"context"

	"wallet-srv/internal/svc"
	"wallet-srv/wallet_srv"

	"github.com/zeromicro/go-zero/core/logx"
)

type WalletsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWalletsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletsLogic {
	return &WalletsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WalletsLogic) Wallets(req *wallet_srv.WalletsReq) (*wallet_srv.WalletsResp, error) {
	walletID, err := l.svcCtx.WalletManager.Init(l.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &wallet_srv.WalletsResp{
		WalletID: walletID,
	}, nil
}
