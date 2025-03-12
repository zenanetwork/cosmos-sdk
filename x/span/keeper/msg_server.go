package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/span/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl은 MsgServer 인터페이스를 구현하는 새로운 인스턴스를 반환합니다.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateSpan은 새로운 스팬을 생성합니다.
func (k msgServer) CreateSpan(goCtx context.Context, msg *types.MsgCreateSpan) (*types.MsgCreateSpanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 권한 검사
	if !k.HasPermission(ctx, msg.Creator) {
		return nil, fmt.Errorf("unauthorized: %s does not have permission to create span", msg.Creator)
	}

	// 현재 스팬 ID 가져오기
	lastSpanID := k.GetLastSpanID(ctx)
	newSpanID := lastSpanID + 1

	// 검증자 세트 생성
	var validators []*types.Validator
	for _, v := range msg.Validators {
		validators = append(validators, &types.Validator{
			Address:          v.Address,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		})
	}

	// 새 스팬 생성
	span := types.NewSpan(
		newSpanID,
		msg.StartBlock,
		msg.EndBlock,
		validators,
		msg.SelectedProducers,
		msg.ChainId,
	)

	// 스팬 저장
	k.SetSpan(ctx, span)
	k.SetLastSpanID(ctx, newSpanID)
	k.SetCurrentSpanID(ctx, newSpanID)

	// 이벤트 발생
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateSpan,
			sdk.NewAttribute(types.AttributeKeySpanID, fmt.Sprintf("%d", newSpanID)),
			sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprintf("%d", msg.StartBlock)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprintf("%d", msg.EndBlock)),
		),
	)

	return &types.MsgCreateSpanResponse{
		Id: newSpanID,
	}, nil
}

// UpdateParams는 모듈 파라미터를 업데이트합니다.
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 권한 검사 - 일반적으로 거버넌스 계정만 파라미터를 업데이트할 수 있음
	if !k.HasPermission(ctx, msg.Authority) {
		return nil, fmt.Errorf("unauthorized: %s does not have permission to update params", msg.Authority)
	}

	// 파라미터 유효성 검사
	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	// 파라미터 업데이트
	k.SetParams(ctx, msg.Params)

	// 이벤트 발생
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateParams,
			sdk.NewAttribute(types.AttributeKeySpanLength, fmt.Sprintf("%d", msg.Params.SpanLength)),
			sdk.NewAttribute(types.AttributeKeyActiveSpanCount, fmt.Sprintf("%d", msg.Params.ActiveSpanCount)),
		),
	)

	return &types.MsgUpdateParamsResponse{}, nil
}

// HasPermission은 주어진 주소가 특정 작업을 수행할 권한이 있는지 확인합니다.
func (k Keeper) HasPermission(ctx sdk.Context, address string) bool {
	// 실제 구현에서는 권한 검사 로직을 구현해야 합니다.
	// 예: 관리자 계정 확인, 거버넌스 제안 확인 등
	// 여기서는 간단히 true를 반환합니다.
	return true
}
