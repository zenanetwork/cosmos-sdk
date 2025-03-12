package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl은 MsgServer 인터페이스를 구현하는 새로운 인스턴스를 반환합니다.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateCheckpoint는 새로운 체크포인트를 생성합니다.
func (k msgServer) CreateCheckpoint(goCtx context.Context, msg *types.MsgCreateCheckpoint) (*types.MsgCreateCheckpointResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 권한 검사
	if !k.HasPermission(ctx, msg.Creator) {
		return nil, fmt.Errorf("unauthorized: %s does not have permission to create checkpoint", msg.Creator)
	}

	// 현재 체크포인트 번호 가져오기
	currentNumber := k.GetCurrentCheckpointNumber(ctx)
	newNumber := currentNumber + 1

	// 새 체크포인트 생성
	checkpoint := types.NewCheckpoint(
		newNumber,
		msg.StartBlock,
		msg.EndBlock,
		msg.RootHash,
		msg.Creator,
	)

	// 체크포인트 저장
	k.SetCheckpoint(ctx, checkpoint)
	k.SetCurrentCheckpointNumber(ctx, newNumber)

	// 이벤트 발생
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCheckpoint,
			sdk.NewAttribute(types.AttributeKeyCheckpointNumber, fmt.Sprintf("%d", newNumber)),
			sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprintf("%d", msg.StartBlock)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprintf("%d", msg.EndBlock)),
		),
	)

	return &types.MsgCreateCheckpointResponse{
		Number: newNumber,
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
			sdk.NewAttribute(types.AttributeKeyCheckpointInterval, fmt.Sprintf("%d", msg.Params.CheckpointInterval)),
			sdk.NewAttribute(types.AttributeKeyCheckpointBufferSize, fmt.Sprintf("%d", msg.Params.CheckpointBufferSize)),
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
