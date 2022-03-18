package handlers

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/service"
	"strings"

	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/dto"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/app/nftp/models/vo"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

type nftTransfer struct {
	base
	svc map[string]service.TransferService
}

type INftTransfer interface {
	TransferNftClassByID(ctx context.Context, _ interface{}) (interface{}, error)
	TransferNftByNftId(ctx context.Context, _ interface{}) (interface{}, error)
}

func NewNftTransfer(svc ...*service.TransferBase) INftTransfer {
	return newNftTransfer(svc)
}

func newNftTransfer(svc []*service.TransferBase) *nftTransfer {
	modules := make(map[string]service.TransferService, len(svc))
	for _, v := range svc {
		modules[v.Module] = v.Service
	}
	return &nftTransfer{
		svc: modules,
	}
}

// TransferNftClassByID transfer an nft class by id
func (h nftTransfer) TransferNftClassByID(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftClassByIDRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipient)
	}
	if len([]rune(recipient)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientLen)
	}
	// 校验接收者地址是否满足当前链的地址规范
	if !common.IsHexAddress(recipient) {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
	}

	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	//校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftClassByIDP{
		ClassID:    h.ClassID(ctx),
		Owner:      h.Owner(ctx),
		Recipient:  recipient,
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Tag:        tagBytes,
	}
	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}
	return service.TransferNFTClass(params)
}

// TransferNftByNftId transfer an nft class by index
func (h nftTransfer) TransferNftByNftId(ctx context.Context, request interface{}) (interface{}, error) {
	// 校验参数 start
	req := request.(*vo.TransferNftByNftIdRequest)
	recipient := strings.TrimSpace(req.Recipient)
	if recipient == "" {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipient)
	}
	if len([]rune(recipient)) > 128 {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientLen)
	}
	// 校验接收者地址是否满足当前链的地址规范
	if !common.IsHexAddress(recipient) {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrRecipientAddr)
	}
	tagBytes, err := h.ValidateTag(req.Tag)
	if err != nil {
		return nil, err
	}

	// 校验参数 end
	authData := h.AuthData(ctx)
	params := dto.TransferNftByNftIdP{
		ClassID:    h.ClassID(ctx),
		Sender:     h.Owner(ctx),
		NftId:      h.NftId(ctx),
		Recipient:  recipient,
		ChainID:    authData.ChainId,
		ProjectID:  authData.ProjectId,
		PlatFormID: authData.PlatformId,
		Tag:        tagBytes,
	}
	//不能自己转让给自己
	//400
	if params.Recipient == params.Sender {
		return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrSelfTransfer)
	}

	service, ok := h.svc[authData.Module]
	if !ok {
		return nil, types.ErrModules
	}

	return service.TransferNFT(params)
}

func (h nftTransfer) ClassID(ctx context.Context) string {
	classId := ctx.Value("class_id")
	if classId == nil {
		return ""
	}
	return classId.(string)
}

func (h nftTransfer) Owner(ctx context.Context) string {
	owner := ctx.Value("owner")
	if owner == nil {
		return ""
	}
	return owner.(string)
}

func (h nftTransfer) NftId(ctx context.Context) string {
	nftId := ctx.Value("nft_id")
	if nftId == nil {
		return ""
	}
	return nftId.(string)
}
