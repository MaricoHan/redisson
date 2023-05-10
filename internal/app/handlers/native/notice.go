package native

import (
	"context"
	"fmt"
	"gitlab.bianjie.ai/avata/open-api/internal/app/handlers"
	"gitlab.bianjie.ai/avata/open-api/internal/app/models/vo/native/notice"
	"gitlab.bianjie.ai/avata/open-api/internal/app/services/native"
	"gitlab.bianjie.ai/avata/utils/errors"
	"strings"
)

type INotice interface {
	TransferNFTS(ctx context.Context, request interface{}) (interface{}, error)
	TransferClasses(ctx context.Context, request interface{}) (interface{}, error)
}

type Notice struct {
	handlers.Base
	handlers.PageBasic
	svc native.INotice
}

func NewNotice(svc native.INotice) *Notice {
	return &Notice{svc: svc}
}

// TransferNFTS 转让NFT通知
func (n Notice) TransferNFTS(ctx context.Context, request interface{}) (interface{}, error) {
	params := request.(*notice.TransferNFTS)
	params.TxHash = strings.TrimSpace(params.TxHash)
	params.ProjectID = strings.TrimSpace(params.ProjectID)
	if params.TxHash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "tx_hash"))
	}
	if params.ProjectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	return n.svc.TransferNFTS(ctx, params)
}

// TransferClasses 转让Class通知
func (n Notice) TransferClasses(ctx context.Context, request interface{}) (interface{}, error) {
	params := request.(*notice.TransferClasses)
	params.TxHash = strings.TrimSpace(params.TxHash)
	params.ProjectID = strings.TrimSpace(params.ProjectID)
	if params.TxHash == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "tx_hash"))
	}
	if params.ProjectID == "" {
		return nil, errors.New(errors.ClientParams, fmt.Sprintf(errors.ErrRequired, "project_id"))
	}
	return n.svc.TransferClasses(ctx, params)
}
