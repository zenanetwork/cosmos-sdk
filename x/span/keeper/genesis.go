package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/span/types"
)

// InitGenesis는 제네시스 상태를 초기화합니다.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	// 파라미터 설정
	k.SetParams(ctx, genState.Params)

	// 스팬 설정
	for _, span := range genState.Spans {
		spanCopy := span // 로컬 변수로 복사
		k.SetSpan(ctx, &spanCopy)
	}

	// 마지막 스팬 ID 설정
	k.SetLastSpanID(ctx, genState.LastSpanID)

	// 현재 스팬 ID 설정
	k.SetCurrentSpanID(ctx, genState.CurrentSpanID)
}

// ExportGenesis는 현재 상태를 제네시스 상태로 내보냅니다.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// 파라미터 가져오기
	params := k.GetParams(ctx)

	// 모든 스팬 가져오기
	spans := k.GetAllSpans(ctx)

	// 마지막 스팬 ID 가져오기
	lastSpanID := k.GetLastSpanID(ctx)

	// 현재 스팬 ID 가져오기
	currentSpanID := k.GetCurrentSpanID(ctx)

	return &types.GenesisState{
		Params:        params,
		Spans:         spans,
		LastSpanID:    lastSpanID,
		CurrentSpanID: currentSpanID,
	}
}
