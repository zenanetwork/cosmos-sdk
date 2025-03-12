package types

import (
	"time"
)

// Checkpoint는 체크포인트 정보를 나타냅니다.
type Checkpoint struct {
	StartBlock uint64    `protobuf:"varint,1,opt,name=start_block,json=startBlock,proto3" json:"start_block,omitempty"`
	EndBlock   uint64    `protobuf:"varint,2,opt,name=end_block,json=endBlock,proto3" json:"end_block,omitempty"`
	RootHash   []byte    `protobuf:"bytes,3,opt,name=root_hash,json=rootHash,proto3" json:"root_hash,omitempty"`
	Proposer   string    `protobuf:"bytes,4,opt,name=proposer,proto3" json:"proposer,omitempty"`
	Timestamp  time.Time `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

// NewCheckpoint는 새로운 Checkpoint 객체를 생성합니다.
func NewCheckpoint(
	startBlock uint64,
	endBlock uint64,
	rootHash []byte,
	proposer string,
) *Checkpoint {
	return &Checkpoint{
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   rootHash,
		Proposer:   proposer,
		Timestamp:  time.Now(),
	}
}
