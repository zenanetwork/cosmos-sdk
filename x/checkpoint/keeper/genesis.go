package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

// InitGenesis는 제네시스 상태를 초기화합니다.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	// 파라미터 설정
	k.SetParams(ctx, genState.Params)

	// 체크포인트 설정
	for _, checkpoint := range genState.Checkpoints {
		k.SetCheckpoint(ctx, &checkpoint)
	}

	// 현재 체크포인트 번호 설정
	k.SetCurrentCheckpointNumber(ctx, genState.CurrentCheckpointNumber)
}

// ExportGenesis는 현재 상태를 제네시스 상태로 내보냅니다.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// 파라미터 가져오기
	params := k.GetParams(ctx)

	// 모든 체크포인트 가져오기
	checkpoints := k.GetAllCheckpoints(ctx)

	// 현재 체크포인트 번호 가져오기
	currentCheckpointNumber := k.GetCurrentCheckpointNumber(ctx)

	return &types.GenesisState{
		Params:                  params,
		Checkpoints:             checkpoints,
		CurrentCheckpointNumber: currentCheckpointNumber,
	}
}
