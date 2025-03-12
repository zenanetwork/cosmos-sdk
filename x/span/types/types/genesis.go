package types

// DefaultGenesis는 기본 제네시스 상태를 반환합니다.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		Spans:      []Span{},
		LastSpanId: 0,
	}
}

// Validate는 제네시스 상태를 검증합니다.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// 스팬 ID 검증
	for _, span := range gs.Spans {
		if span.Id > gs.LastSpanId {
			return ErrInvalidSpanID
		}
	}

	return nil
}

// NewGenesisState는 새로운 제네시스 상태를 생성합니다.
func NewGenesisState(params Params, spans []Span, lastSpanID uint64) *GenesisState {
	return &GenesisState{
		Params:     params,
		Spans:      spans,
		LastSpanId: lastSpanID,
	}
}
