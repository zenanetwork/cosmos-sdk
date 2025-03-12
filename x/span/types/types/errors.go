package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// 모듈별 오류 코드 정의
var (
	// 1-100 범위는 모듈 자체 오류용으로 예약됨
	ErrInvalidSpanID       = sdkerrors.Register(ModuleName, 1, "invalid span ID")
	ErrSpanNotFound        = sdkerrors.Register(ModuleName, 2, "span not found")
	ErrInvalidSpanLength   = sdkerrors.Register(ModuleName, 3, "invalid span length")
	ErrInvalidBlockRange   = sdkerrors.Register(ModuleName, 4, "invalid block range")
	ErrInvalidValidatorSet = sdkerrors.Register(ModuleName, 5, "invalid validator set")
	ErrInvalidProducer     = sdkerrors.Register(ModuleName, 6, "invalid producer")
	ErrUnauthorized        = sdkerrors.Register(ModuleName, 7, "unauthorized")
)
