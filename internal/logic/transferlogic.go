package logic

import (
	"context"
	"errors"

	"wallet-srv/internal/dao"
	"wallet-srv/internal/svc"
	"wallet-srv/wallet_srv"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransferLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferLogic {
	return &TransferLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransferLogic) Transfer(req *wallet_srv.TransferReq) (*wallet_srv.TransferResp, error) {
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "amount 格式不合法")
	}
	err = l.svcCtx.WalletManager.Transfer(l.ctx, req.FromWalletID, req.ToWalletID, amount)
	if err != nil {
		switch {
		case errors.Is(err, dao.ErrWalletNotFound):
			return nil, status.Error(codes.NotFound, "钱包不存在")
		case errors.Is(err, dao.ErrInsufficientBalance):
			return nil, status.Error(codes.FailedPrecondition, "余额不足")
		case errors.Is(err, dao.ErrInvalidAmount), errors.Is(err, dao.ErrInvalidWalletID):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, dao.ErrSameWalletTransfer):
			return nil, status.Error(codes.InvalidArgument, "from/to 不能相同")
		default:
			return nil, status.Error(codes.Internal, "转账失败")
		}
	}
	return &wallet_srv.TransferResp{
		Result: "ok",
	}, nil
}
