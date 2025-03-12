package types

import (
	"encoding/binary"
	"fmt"
)

const (
	// ModuleName은 모듈의 이름입니다.
	ModuleName = "span"

	// StoreKey는 상태 저장소에 접근하기 위한 키입니다.
	StoreKey = ModuleName

	// RouterKey는 메시지 라우팅을 위한 키입니다.
	RouterKey = ModuleName

	// QuerierRoute는 쿼리 라우팅을 위한 키입니다.
	QuerierRoute = ModuleName
)

var (
	// SpanKeyPrefix는 스팬 정보를 저장하기 위한 키 접두사입니다.
	SpanKeyPrefix = []byte{0x01}

	// SpanCountKey는 스팬 수를 저장하기 위한 키입니다.
	SpanCountKey = []byte{0x02}

	// LastSpanIDKey는 마지막 스팬 ID를 저장하기 위한 키입니다.
	LastSpanIDKey = []byte{0x03}
)

// GetSpanKey는 주어진 스팬 ID에 대한 키를 반환합니다.
func GetSpanKey(spanID uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, spanID)
	return append(SpanKeyPrefix, bz...)
}

// GetSpanIDFromKey는 키에서 스팬 ID를 추출합니다.
func GetSpanIDFromKey(key []byte) (uint64, error) {
	if len(key) != len(SpanKeyPrefix)+8 {
		return 0, fmt.Errorf("invalid key length: %d", len(key))
	}
	return binary.BigEndian.Uint64(key[len(SpanKeyPrefix):]), nil
}
