package keeper

import (
	"encoding/binary"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

// Keeper는 checkpoint 모듈의 상태를 관리합니다.
type Keeper struct {
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	authority string

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	spanKeeper    types.SpanKeeper
}

// NewKeeper는 새로운 Keeper 인스턴스를 생성합니다.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	spanKeeper types.SpanKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		authority:     authority,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		spanKeeper:    spanKeeper,
	}
}

// GetAuthority는 모듈 권한 주소를 반환합니다.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// GetParams는 모듈 파라미터를 반환합니다.
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
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
}

// GetCheckpoint는 주어진 번호에 해당하는 체크포인트를 반환합니다.
func (k Keeper) GetCheckpoint(ctx sdk.Context, number int64) (*types.Checkpoint, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.CheckpointKey(number)
	bz := store.Get(key)
	if bz == nil {
		return nil, fmt.Errorf("checkpoint with number %d not found", number)
	}

	var checkpoint types.Checkpoint
	k.cdc.MustUnmarshal(bz, &checkpoint)
	return &checkpoint, nil
}

// SetCheckpoint는 체크포인트를 저장합니다.
func (k Keeper) SetCheckpoint(ctx sdk.Context, checkpoint *types.Checkpoint) {
	store := ctx.KVStore(k.storeKey)
	key := types.CheckpointKey(checkpoint.Number)
	bz := k.cdc.MustMarshal(checkpoint)
	store.Set(key, bz)
}

// GetAllCheckpoints는 모든 체크포인트를 반환합니다.
func (k Keeper) GetAllCheckpoints(ctx sdk.Context) []types.Checkpoint {
	var checkpoints []types.Checkpoint
	store := ctx.KVStore(k.storeKey)

	// 직접 이터레이터 구현
	iterator := store.Iterator(types.CheckpointKeyPrefix, storetypes.PrefixEndBytes(types.CheckpointKeyPrefix))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var checkpoint types.Checkpoint
		k.cdc.MustUnmarshal(iterator.Value(), &checkpoint)
		checkpoints = append(checkpoints, checkpoint)
	}

	return checkpoints
}

// GetCurrentCheckpointNumber는 현재 체크포인트 번호를 반환합니다.
func (k Keeper) GetCurrentCheckpointNumber(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrentCheckpointNumberKey)
	if bz == nil {
		return 0
	}

	return int64(binary.BigEndian.Uint64(bz))
}

// SetCurrentCheckpointNumber는 현재 체크포인트 번호를 설정합니다.
func (k Keeper) SetCurrentCheckpointNumber(ctx sdk.Context, number int64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(number))
	store.Set(types.CurrentCheckpointNumberKey, bz)
}

// GetLatestCheckpoint는 최신 체크포인트를 가져옵니다.
func (k Keeper) GetLatestCheckpoint(ctx sdk.Context) (*types.Checkpoint, error) {
	currentNumber := k.GetCurrentCheckpointNumber(ctx)
	return k.GetCheckpoint(ctx, currentNumber)
}

// CreateCheckpoint는 새로운 체크포인트를 생성하고 저장합니다.
func (k Keeper) CreateCheckpoint(
	ctx sdk.Context,
	number int64,
	startBlock uint64,
	endBlock uint64,
	rootHash []byte,
	proposer string,
) *types.Checkpoint {
	checkpoint := types.NewCheckpoint(
		number,
		startBlock,
		endBlock,
		rootHash,
		proposer,
	)

	k.SetCheckpoint(ctx, checkpoint)

	// 이벤트 발행
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateCheckpoint,
			sdk.NewAttribute(types.AttributeKeyCheckpointNumber, fmt.Sprintf("%d", checkpoint.Number)),
			sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprintf("%d", checkpoint.StartBlock)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprintf("%d", checkpoint.EndBlock)),
			sdk.NewAttribute(types.AttributeKeyProposer, checkpoint.Proposer),
		),
	)

	return checkpoint
}

// GetCheckpointCount는 체크포인트 수를 가져옵니다.
func (k Keeper) GetCheckpointCount(ctx sdk.Context) int64 {
	return k.GetCurrentCheckpointNumber(ctx) + 1
}
