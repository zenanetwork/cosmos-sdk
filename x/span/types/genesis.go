package types

import (
	"github.com/cosmos/gogoproto/proto"
)

// GenesisState는 span 모듈의 제네시스 상태를 정의합니다.
type GenesisState struct {
	Params        Params `json:"params" yaml:"params"`
	Spans         []Span `json:"spans" yaml:"spans"`
	LastSpanID    uint64 `json:"last_span_id" yaml:"last_span_id"`
	CurrentSpanID uint64 `json:"current_span_id"`
}

// ProtoMessage 인터페이스 구현
func (*GenesisState) ProtoMessage()    {}
func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }

// DefaultGenesis는 기본 제네시스 상태를 반환합니다.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:        DefaultParams(),
		Spans:         []Span{},
		LastSpanID:    0,
		CurrentSpanID: 0,
	}
}

// Validate는 제네시스 상태를 검증합니다.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// 스팬 ID 검증
	for _, span := range gs.Spans {
		if span.Id > gs.LastSpanID {
			return ErrInvalidSpanID
		}

		if span.StartBlock >= span.EndBlock {
			return ErrInvalidBlockRange
		}

		if len(span.ValidatorSet) == 0 {
			return ErrInvalidValidatorSet
		}
	}

	return nil
}

// NewGenesisState는 새로운 제네시스 상태를 생성합니다.
func NewGenesisState(params Params, spans []Span, lastSpanID uint64) *GenesisState {
	return &GenesisState{
		Params:        params,
		Spans:         spans,
		LastSpanID:    lastSpanID,
		CurrentSpanID: 0,
	}
}
