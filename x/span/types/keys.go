package types

import (
	"encoding/binary"
)

const (
	// ModuleName은 모듈의 이름입니다.
	ModuleName = "span"

	// StoreKey는 모듈의 스토어 키입니다.
	StoreKey = ModuleName

	// RouterKey는 모듈의 라우터 키입니다.
	RouterKey = ModuleName

	// QuerierRoute는 모듈의 쿼리 라우트입니다.
	QuerierRoute = ModuleName
)

var (
	// SpanKeyPrefix는 스팬 키의 접두사입니다.
	SpanKeyPrefix = []byte{0x01}

	// SpanCountKey는 스팬 수를 저장하는 키입니다.
	SpanCountKey = []byte{0x02}

	// LastSpanIDKey는 마지막 스팬 ID를 저장하는 키입니다.
	LastSpanIDKey = []byte{0x03}

	// CurrentSpanIDKey는 현재 스팬 ID를 저장하는 키입니다.
	CurrentSpanIDKey = []byte{0x04}

	// ParamsKey는 모듈 파라미터를 저장하는 키입니다.
	ParamsKey = []byte{0x05}
)

// SpanKey는 주어진 ID에 대한 스팬 키를 반환합니다.
func SpanKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append(SpanKeyPrefix, bz...)
}

// GetSpanIDFromKey는 키에서 스팬 ID를 추출합니다.
func GetSpanIDFromKey(key []byte) (uint64, error) {
	if len(key) != 9 {
		return 0, ErrInvalidKey
	}
	return binary.BigEndian.Uint64(key[1:]), nil
}
