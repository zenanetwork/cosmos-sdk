package types

import (
	"errors"
)

// 모듈별 오류 정의
var (
	// 기본 오류
	ErrInvalidBlockRange  = errors.New("invalid block range")
	ErrInvalidRootHash    = errors.New("invalid root hash")
	ErrInvalidKey         = errors.New("invalid key")
	ErrCheckpointNotFound = errors.New("checkpoint not found")
	ErrUnauthorized       = errors.New("unauthorized")
)
