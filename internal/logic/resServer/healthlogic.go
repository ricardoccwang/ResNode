package resServer

import (
	"context"

	"UResNode/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
)

type HealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) HealthLogic {
	return HealthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthLogic) Health() error {
	l.svcCtx.Client.Cure()
	return nil
}
