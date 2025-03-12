package types

import (
	"fmt"
	"time"
)

// Checkpoint는 블록 체인의 특정 지점에서의 상태 스냅샷을 나타냅니다.
type Checkpoint struct {
	Number     int64     `json:"number"`
	StartBlock uint64    `json:"start_block"`
	EndBlock   uint64    `json:"end_block"`
	RootHash   []byte    `json:"root_hash"`
	Proposer   string    `json:"proposer"`
	Timestamp  time.Time `json:"timestamp"`
}

// ProtoMessage 인터페이스 구현
func (m *Checkpoint) ProtoMessage() {}
func (m *Checkpoint) Reset()        { *m = Checkpoint{} }
func (m *Checkpoint) String() string {
	return fmt.Sprintf("Checkpoint{Number: %d, StartBlock: %d, EndBlock: %d}",
		m.Number, m.StartBlock, m.EndBlock)
}

// NewCheckpoint는 새로운 Checkpoint 객체를 생성합니다.
func NewCheckpoint(
	number int64,
	startBlock uint64,
	endBlock uint64,
	rootHash []byte,
	proposer string,
) *Checkpoint {
	return &Checkpoint{
		Number:     number,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   rootHash,
		Proposer:   proposer,
		Timestamp:  time.Now(),
	}
}
