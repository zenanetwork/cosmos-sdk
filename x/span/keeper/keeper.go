package keeper

import (
	"encoding/binary"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/zenanetwork/cosmos-sdk/codec"
	sdk "github.com/zenanetwork/cosmos-sdk/types"
	"github.com/zenanetwork/cosmos-sdk/x/span/types"
)

// Keeper는 span 모듈의 상태를 관리합니다.
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	authority     string
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewKeeper는 새로운 Keeper를 생성합니다.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		authority:     authority,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// GetAuthority는 모듈 권한 주소를 반환합니다.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams는 모듈 파라미터를 가져옵니다.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams()
	}

	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams는 모듈 파라미터를 설정합니다.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetCurrentSpanID는 현재 스팬 ID를 가져옵니다.
func (k Keeper) GetCurrentSpanID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrentSpanIDKey)
	if bz == nil {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// SetCurrentSpanID는 현재 스팬 ID를 설정합니다.
func (k Keeper) SetCurrentSpanID(ctx sdk.Context, spanID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, spanID)
	store.Set(types.CurrentSpanIDKey, bz)
}

// GetSpan은 특정 ID의 스팬을 가져옵니다.
func (k Keeper) GetSpan(ctx sdk.Context, spanID uint64) (*types.Span, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SpanKey(spanID))
	if bz == nil {
		return nil, false
	}

	var span types.Span
	k.cdc.MustUnmarshal(bz, &span)
	return &span, true
}

// SetSpan은 스팬을 저장합니다.
func (k Keeper) SetSpan(ctx sdk.Context, span types.Span) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&span)
	store.Set(types.SpanKey(span.Id), bz)
}

// GetCurrentSpan은 현재 활성화된 스팬을 가져옵니다.
func (k Keeper) GetCurrentSpan(ctx sdk.Context) (*types.Span, bool) {
	currentSpanID := k.GetCurrentSpanID(ctx)
	return k.GetSpan(ctx, currentSpanID)
}

// GetSpanByHeight는 주어진 높이에 해당하는 스팬을 반환합니다.
func (k Keeper) GetSpanByHeight(ctx sdk.Context, height uint64) (*types.Span, bool) {
	currentSpanID := k.GetCurrentSpanID(ctx)
	if currentSpanID == 0 {
		return nil, false
	}

	// 현재 스팬부터 역순으로 검색
	for id := currentSpanID; id > 0; id-- {
		span, found := k.GetSpan(ctx, id)
		if !found {
			continue
		}

		// 높이가 스팬 범위 내에 있는지 확인
		if height >= span.StartBlock && height <= span.EndBlock {
			return span, true
		}

		// 높이가 현재 스팬보다 크면 더 이상 검색할 필요 없음
		if height > span.EndBlock {
			break
		}
	}

	return nil, false
}

// CreateSpan은 새로운 스팬을 생성하고 저장합니다.
func (k Keeper) CreateSpan(
	ctx sdk.Context,
	id uint64,
	startBlock uint64,
	endBlock uint64,
	validatorSet []*types.Validator,
	selectedProducers []string,
	chainID string,
) *types.Span {
	// 체인 ID가 비어있으면 파라미터에서 가져옴
	if chainID == "" {
		params := k.GetParams(ctx)
		chainID = params.ChainID
	}

	span := types.NewSpan(
		id,
		startBlock,
		endBlock,
		validatorSet,
		selectedProducers,
		chainID,
	)

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(span)
	store.Set(types.SpanKey(id), bz)

	// 이벤트 발행
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateSpan,
			sdk.NewAttribute(types.AttributeKeySpanID, fmt.Sprintf("%d", span.Id)),
			sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprintf("%d", span.StartBlock)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprintf("%d", span.EndBlock)),
			sdk.NewAttribute(types.AttributeKeyChainID, span.ChainId),
		),
	)

	return span
}

// GetAllSpans는 모든 스팬을 가져옵니다.
func (k Keeper) GetAllSpans(ctx sdk.Context) []types.Span {
	store := ctx.KVStore(k.storeKey)

	// 직접 이터레이터 구현
	iterator := store.Iterator(types.SpanKeyPrefix, storetypes.PrefixEndBytes(types.SpanKeyPrefix))
	defer iterator.Close()

	spans := []types.Span{}
	for ; iterator.Valid(); iterator.Next() {
		var span types.Span
		k.cdc.MustUnmarshal(iterator.Value(), &span)
		spans = append(spans, span)
	}

	return spans
}

// GetLastSpanID는 마지막 스팬 ID를 반환합니다.
func (k Keeper) GetLastSpanID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastSpanIDKey)
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// SetLastSpanID는 마지막 스팬 ID를 설정합니다.
func (k Keeper) SetLastSpanID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	store.Set(types.LastSpanIDKey, bz)
}
