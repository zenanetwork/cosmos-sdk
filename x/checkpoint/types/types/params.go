package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// 기본 파라미터 값
const (
	DefaultCheckpointInterval = uint64(1024) // 기본 체크포인트 간격
)

// 파라미터 키
var (
	KeyCheckpointInterval = []byte("CheckpointInterval")
)

// ParamKeyTable는 파라미터 키 테이블을 반환합니다.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams는 새로운 파라미터 객체를 생성합니다.
func NewParams(checkpointInterval uint64) Params {
	return Params{
		CheckpointInterval: checkpointInterval,
	}
}

// DefaultParams는 기본 파라미터 값을 반환합니다.
func DefaultParams() Params {
	return NewParams(DefaultCheckpointInterval)
}

// ParamSetPairs는 파라미터 세트 쌍을 구현합니다.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCheckpointInterval, &p.CheckpointInterval, validateCheckpointInterval),
	}
}

// Validate는 파라미터 값을 검증합니다.
func (p Params) Validate() error {
	if err := validateCheckpointInterval(p.CheckpointInterval); err != nil {
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

// Params는 checkpoint 모듈의 파라미터를 정의합니다.
type Params struct {
	CheckpointInterval uint64 `protobuf:"varint,1,opt,name=checkpoint_interval,json=checkpointInterval,proto3" json:"checkpoint_interval,omitempty"`
}
