package types

import (
	"errors"
)

// 모듈별 오류 정의
var (
	// 기본 오류
	ErrInvalidSpanID       = errors.New("invalid span ID")
	ErrInvalidBlockRange   = errors.New("invalid block range")
	ErrInvalidValidatorSet = errors.New("invalid validator set")
	ErrInvalidKey          = errors.New("invalid key")
	ErrSpanNotFound        = errors.New("span not found")
	ErrUnauthorized        = errors.New("unauthorized")
)
