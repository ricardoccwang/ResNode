package logic

import (
	"context"

	"UResNode/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
)

type TestRefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) TestRefreshTokenLogic {
	return TestRefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestRefreshTokenLogic) TestRefreshToken() error {
	l.svcCtx.Client.Token.ReFreshToken()

	return nil
}
