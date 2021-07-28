package logic

import (
	"context"

	"github.com/aaronlyc/allocation/server/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
)

type CleanAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCleanAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) CleanAllLogic {
	return CleanAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanAllLogic) CleanAll() error {
	// todo: add your logic here and delete this line

	return nil
}
