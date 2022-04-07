package resServer

import (
	"UResNode/UStorageClient"
	"UResNode/UTool"
	"UResNode/internal/Data"
	"context"
	"net/http"

	"UResNode/internal/svc"
	"UResNode/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UploadResLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadResLogic(ctx context.Context, svcCtx *svc.ServiceContext) UploadResLogic {
	return UploadResLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadResLogic) Fail(message string) *types.UploadResResp {
	return &types.UploadResResp{
		GeneralResponse: types.GeneralResponse{
			Code:    Data.API_FAIL,
			Message: message,
		},
	}
}

func (l *UploadResLogic) Success(res *UStorageClient.ResObject) *types.UploadResResp {
	return &types.UploadResResp{
		GeneralResponse: types.GeneralResponse{
			Code:    Data.API_SUCCESS,
			Message: "",
		},
		Name:          res.Name,
		FileUrl:       res.FileUrl,
		FileSignature: res.FileSignature,
		FileSize:      res.FileSize,
	}
}

func (l *UploadResLogic) UploadRes(req types.UploadResReq, r *http.Request) (*types.UploadResResp, error) {
	if ok, unit, err := l.svcCtx.Client.GetResUnitBySignature(req.Signature); err != nil {
		return l.Fail(err.Error()), nil
	} else if ok {
		return l.Success(UStorageClient.CreateResObjectFromResUnit(l.svcCtx.Client.DownloadUrl, unit)), nil
	} else {
		saveFilePath, err := l.svcCtx.Client.SaveFileFromRequest(req, r)
		if err != nil {
			return l.Fail(UTool.LogxBothFail(err.Error()).Error()), nil
		}
		res, err := l.svcCtx.Client.AddNewFileIndex(saveFilePath)
		if err != nil {
			return l.Fail(UTool.LogxBothFail(err.Error()).Error()), nil
		}
		return l.Success(UStorageClient.CreateResObjectFromResUnit(l.svcCtx.Client.DownloadUrl, res)), nil
	}
}
