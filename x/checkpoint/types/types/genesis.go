package types

// DefaultGenesis는 기본 제네시스 상태를 반환합니다.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		Checkpoints:          []Checkpoint{},
		LastCheckpointNumber: 0,
	}
}

// Validate는 제네시스 상태를 검증합니다.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// 체크포인트 번호 검증
	for _, checkpoint := range gs.Checkpoints {
		if checkpoint.StartBlock > checkpoint.EndBlock {
			return ErrInvalidBlockRange
		}
	}

	return nil
}

// NewGenesisState는 새로운 제네시스 상태를 생성합니다.
func NewGenesisState(params Params, checkpoints []Checkpoint, lastCheckpointNumber int64) *GenesisState {
	return &GenesisState{
		Params:               params,
		Checkpoints:          checkpoints,
		LastCheckpointNumber: lastCheckpointNumber,
	}
}
