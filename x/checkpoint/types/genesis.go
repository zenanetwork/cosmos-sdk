package types

import (
	"github.com/cosmos/gogoproto/proto"
)

// GenesisState는 모듈의 제네시스 상태를 정의합니다.
type GenesisState struct {
	Params                  Params       `json:"params" yaml:"params"`
	Checkpoints             []Checkpoint `json:"checkpoints" yaml:"checkpoints"`
	CurrentCheckpointNumber int64        `json:"current_checkpoint_number" yaml:"current_checkpoint_number"`
}

// ProtoMessage 인터페이스 구현
func (*GenesisState) ProtoMessage()    {}
func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }

// DefaultGenesis는 기본 제네시스 상태를 반환합니다.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                  DefaultParams(),
		Checkpoints:             []Checkpoint{},
		CurrentCheckpointNumber: 0,
	}
}

// Validate는 제네시스 상태의 유효성을 검사합니다.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// 체크포인트 유효성 검사
	for _, checkpoint := range gs.Checkpoints {
		if checkpoint.StartBlock >= checkpoint.EndBlock {
			return ErrInvalidBlockRange
		}

		if len(checkpoint.RootHash) == 0 {
			return ErrInvalidRootHash
		}
	}

	return nil
}

// NewGenesisState는 새로운 제네시스 상태를 생성합니다.
func NewGenesisState(params Params, checkpoints []Checkpoint, lastCheckpointNumber int64) *GenesisState {
	return &GenesisState{
		Params:                  params,
		Checkpoints:             checkpoints,
		CurrentCheckpointNumber: lastCheckpointNumber,
	}
}
