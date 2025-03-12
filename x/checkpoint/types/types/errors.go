package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// 모듈별 오류 코드 정의
var (
	// 1-100 범위는 모듈 자체 오류용으로 예약됨
	ErrInvalidCheckpointNumber   = sdkerrors.Register(ModuleName, 1, "invalid checkpoint number")
	ErrCheckpointNotFound        = sdkerrors.Register(ModuleName, 2, "checkpoint not found")
	ErrInvalidCheckpointInterval = sdkerrors.Register(ModuleName, 3, "invalid checkpoint interval")
	ErrInvalidBlockRange         = sdkerrors.Register(ModuleName, 4, "invalid block range")
	ErrInvalidRootHash           = sdkerrors.Register(ModuleName, 5, "invalid root hash")
	ErrUnauthorized              = sdkerrors.Register(ModuleName, 6, "unauthorized")
)
