package types

import (
	"fmt"
)

// 기본 파라미터 값
const (
	DefaultCheckpointInterval = uint64(1024) // 기본 체크포인트 간격
	DefaultChainID            = "zenachain"  // 기본 체인 ID
)

// 파라미터 키
var (
	KeyCheckpointInterval   = []byte("CheckpointInterval")
	KeyCheckpointBufferSize = []byte("CheckpointBufferSize")
	KeyChainID              = []byte("ChainID")
)

// Params는 checkpoint 모듈의 파라미터를 정의합니다.
type Params struct {
	CheckpointInterval   uint64 `json:"checkpoint_interval"`
	CheckpointBufferSize uint64 `json:"checkpoint_buffer_size"`
	ChainID              string `json:"chain_id"`
}

// ProtoMessage 인터페이스 구현
func (m *Params) ProtoMessage() {}
func (m *Params) Reset()        { *m = Params{} }
func (m *Params) String() string {
	return fmt.Sprintf("Params{CheckpointInterval: %d, CheckpointBufferSize: %d, ChainID: %s}",
		m.CheckpointInterval, m.CheckpointBufferSize, m.ChainID)
}

// NewParams는 새로운 파라미터 객체를 생성합니다.
func NewParams(checkpointInterval uint64, checkpointBufferSize uint64, chainID string) Params {
	return Params{
		CheckpointInterval:   checkpointInterval,
		CheckpointBufferSize: checkpointBufferSize,
		ChainID:              chainID,
	}
}

// DefaultParams는 기본 파라미터 값을 반환합니다.
func DefaultParams() Params {
	return Params{
		CheckpointInterval:   100,
		CheckpointBufferSize: 10,
		ChainID:              DefaultChainID,
	}
}

// Validate는 파라미터 값을 검증합니다.
func (p Params) Validate() error {
	if err := validateCheckpointInterval(p.CheckpointInterval); err != nil {
		return err
	}

	if err := validateCheckpointBufferSize(p.CheckpointBufferSize); err != nil {
		return err
	}

	if err := validateChainID(p.ChainID); err != nil {
		return err
	}

	return nil
}

// validateCheckpointInterval는 체크포인트 간격 값을 검증합니다.
func validateCheckpointInterval(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("checkpoint interval must be positive: %d", v)
	}

	return nil
}

func validateCheckpointBufferSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("checkpoint buffer size must be positive: %d", v)
	}

	return nil
}

// validateChainID는 체인 ID의 유효성을 검사합니다.
func validateChainID(chainID string) error {
	if chainID == "" {
		return fmt.Errorf("chain ID cannot be empty")
	}
	return nil
}
