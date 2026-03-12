package svc

import (
	"wallet-srv/internal/config"
	"wallet-srv/internal/dao"
)

type ServiceContext struct {
	Config        config.Config
	WalletManager dao.IWallet
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		WalletManager: dao.NewWalletManager(),
	}
}
