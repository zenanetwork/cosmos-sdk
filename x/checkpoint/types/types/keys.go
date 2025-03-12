package types

import (
	"encoding/binary"
	"fmt"
)

const (
	// ModuleName은 모듈의 이름입니다.
	ModuleName = "checkpoint"

	// StoreKey는 상태 저장소에 접근하기 위한 키입니다.
	StoreKey = ModuleName

	// RouterKey는 메시지 라우팅을 위한 키입니다.
	RouterKey = ModuleName

	// QuerierRoute는 쿼리 라우팅을 위한 키입니다.
	QuerierRoute = ModuleName
)

var (
	// CheckpointKeyPrefix는 체크포인트 정보를 저장하기 위한 키 접두사입니다.
	CheckpointKeyPrefix = []byte{0x01}

	// CheckpointCountKey는 체크포인트 수를 저장하기 위한 키입니다.
	CheckpointCountKey = []byte{0x02}

	// LastCheckpointKey는 마지막 체크포인트를 저장하기 위한 키입니다.
	LastCheckpointKey = []byte{0x03}
)

// GetCheckpointKey는 주어진 체크포인트 번호에 대한 키를 반환합니다.
func GetCheckpointKey(number int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(number))
	return append(CheckpointKeyPrefix, bz...)
}

// GetCheckpointNumberFromKey는 키에서 체크포인트 번호를 추출합니다.
func GetCheckpointNumberFromKey(key []byte) (int64, error) {
	if len(key) != len(CheckpointKeyPrefix)+8 {
		return 0, fmt.Errorf("invalid key length: %d", len(key))
	}
	return int64(binary.BigEndian.Uint64(key[len(CheckpointKeyPrefix):])), nil
}
