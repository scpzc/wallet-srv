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

func (l *WalletsLogic) Wallets(in *wallet_srv.WalletsReq) (*wallet_srv.WalletsResp, error) {
	// todo: add your logic here and delete this line

	return &wallet_srv.WalletsResp{}, nil
}
