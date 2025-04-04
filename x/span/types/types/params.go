package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// 기본 파라미터 값
const (
	DefaultSpanLength = uint64(100) // 기본 스팬 길이
)

// 파라미터 키
var (
	KeySpanLength = []byte("SpanLength")
)

// ParamKeyTable는 파라미터 키 테이블을 반환합니다.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams는 새로운 파라미터 객체를 생성합니다.
func NewParams(spanLength uint64) Params {
	return Params{
		SpanLength: spanLength,
	}
}

// DefaultParams는 기본 파라미터 값을 반환합니다.
func DefaultParams() Params {
	return NewParams(DefaultSpanLength)
}

// ParamSetPairs는 파라미터 세트 쌍을 구현합니다.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySpanLength, &p.SpanLength, validateSpanLength),
	}
}

// Validate는 파라미터 값을 검증합니다.
func (p Params) Validate() error {
	if err := validateSpanLength(p.SpanLength); err != nil {
		return err
	}
	return nil
}

// validateSpanLength는 스팬 길이 값을 검증합니다.
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

// Params는 span 모듈의 파라미터를 정의합니다.
type Params struct {
	SpanLength uint64 `protobuf:"varint,1,opt,name=span_length,json=spanLength,proto3" json:"span_length,omitempty"`
}
