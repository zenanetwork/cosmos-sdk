package types

import (
	"fmt"
)

// 기본 파라미터 값
const (
	DefaultSpanLength = uint64(100) // 기본 스팬 길이
	DefaultChainID    = "zenachain" // 기본 체인 ID
)

// 파라미터 스토어 키
var (
	KeySpanLength      = []byte("SpanLength")
	KeyActiveSpanCount = []byte("ActiveSpanCount")
	KeyChainID         = []byte("ChainID")
)

// Params는 모듈 파라미터를 정의합니다.
type Params struct {
	SpanLength      uint64 `json:"span_length"`
	ActiveSpanCount uint64 `json:"active_span_count"`
	ChainID         string `json:"chain_id"`
}

// ProtoMessage 인터페이스 구현
func (m *Params) ProtoMessage() {}
func (m *Params) Reset()        { *m = Params{} }
func (m *Params) String() string {
	return fmt.Sprintf("Params{SpanLength: %d, ActiveSpanCount: %d, ChainID: %s}",
		m.SpanLength, m.ActiveSpanCount, m.ChainID)
}

// DefaultParams는 기본 파라미터를 반환합니다.
func DefaultParams() Params {
	return Params{
		SpanLength:      100,
		ActiveSpanCount: 10,
		ChainID:         DefaultChainID,
	}
}

// Validate는 파라미터의 유효성을 검사합니다.
func (p Params) Validate() error {
	if err := validateSpanLength(p.SpanLength); err != nil {
		return err
	}

	if err := validateActiveSpanCount(p.ActiveSpanCount); err != nil {
		return err
	}

	if err := validateChainID(p.ChainID); err != nil {
		return err
	}

	return nil
}

func validateSpanLength(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("span length must be positive: %d", v)
	}

	return nil
}

func validateActiveSpanCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("active span count must be positive: %d", v)
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
