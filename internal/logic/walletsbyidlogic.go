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
		switch {
		case errors.Is(err, dao.ErrWalletNotFound):
			return nil, status.Error(codes.NotFound, "钱包不存在")
		case errors.Is(err, dao.ErrInvalidWalletID):
			return nil, status.Error(codes.InvalidArgument, "wallet_id 不合法")
		default:
			return nil, status.Error(codes.Internal, "查询失败")
		}
	}
	return &wallet_srv.WalletsByIDResp{
		WalletID: wallet.UserId,
		Balance:  wallet.Balance.String(),
	}, nil
}
