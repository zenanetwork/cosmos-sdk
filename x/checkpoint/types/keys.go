package types

import (
	"encoding/binary"
)

const (
	// ModuleName은 모듈의 이름입니다.
	ModuleName = "checkpoint"

	// StoreKey는 모듈의 스토어 키입니다.
	StoreKey = ModuleName

	// RouterKey는 모듈의 라우터 키입니다.
	RouterKey = ModuleName

	// QuerierRoute는 모듈의 쿼리 라우트입니다.
	QuerierRoute = ModuleName
)

var (
	// CheckpointKeyPrefix는 체크포인트 키의 접두사입니다.
	CheckpointKeyPrefix = []byte{0x01}

	// CheckpointCountKey는 체크포인트 수를 저장하는 키입니다.
	CheckpointCountKey = []byte{0x02}

	// CurrentCheckpointNumberKey는 현재 체크포인트 번호를 저장하는 키입니다.
	CurrentCheckpointNumberKey = []byte{0x03}

	// ParamsKey는 모듈 파라미터를 저장하는 키입니다.
	ParamsKey = []byte{0x04}
)

// CheckpointKey는 주어진 번호에 대한 체크포인트 키를 반환합니다.
func CheckpointKey(number int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(number))
	return append(CheckpointKeyPrefix, bz...)
}

// GetCheckpointNumberFromKey는 키에서 체크포인트 번호를 추출합니다.
func GetCheckpointNumberFromKey(key []byte) (int64, error) {
	if len(key) != 9 {
		return 0, ErrInvalidKey
	}
	return int64(binary.BigEndian.Uint64(key[1:])), nil
}

// Int64ToBytes는 int64 값을 바이트 슬라이스로 변환합니다.
func Int64ToBytes(val int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(val))
	return bz
}
