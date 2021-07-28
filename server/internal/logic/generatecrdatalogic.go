package logic

import (
	"context"

	"github.com/aaronlyc/allocation/server/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
)

type GenerateCRDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCRDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) GenerateCRDataLogic {
	return GenerateCRDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCRDataLogic) GenerateCRData() error {
	// todo: add your logic here and delete this line

	return nil
}
