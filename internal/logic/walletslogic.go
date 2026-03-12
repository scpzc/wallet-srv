package logic

import (
	"context"
	"errors"

	"wallet-srv/internal/dao"
	"wallet-srv/internal/svc"
	"wallet-srv/wallet_srv"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		switch {
		case errors.Is(err, dao.ErrInvalidWalletID):
			return nil, status.Error(codes.InvalidArgument, "user_id 不合法")
		default:
			return nil, status.Error(codes.Internal, "初始化钱包失败")
		}
	}
	return &wallet_srv.WalletsResp{
		WalletID: walletID,
	}, nil
}
