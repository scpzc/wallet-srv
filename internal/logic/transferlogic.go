package logic

import (
	"context"

	"wallet-srv/internal/svc"
	"wallet-srv/wallet_srv"

	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, err
	}
	err = l.svcCtx.WalletManager.Transfer(l.ctx, req.FromWalletID, req.ToWalletID, amount)
	if err != nil {
		return nil, err
	}
	return &wallet_srv.TransferResp{
		Result: "ok",
	}, nil
}
