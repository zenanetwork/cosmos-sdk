package types

import (
	"fmt"
	"time"
)

// Validator는 검증자 정보를 나타냅니다.
type Validator struct {
	Address          string `json:"address"`
	VotingPower      int64  `json:"voting_power"`
	ProposerPriority int64  `json:"proposer_priority"`
}

// ProtoMessage 인터페이스 구현
func (m *Validator) ProtoMessage() {}
func (m *Validator) Reset()        { *m = Validator{} }
func (m *Validator) String() string {
	return fmt.Sprintf("Validator{Address: %s, VotingPower: %d}",
		m.Address, m.VotingPower)
}

// Span은 블록 범위와 관련 정보를 나타냅니다.
type Span struct {
	Id                uint64       `json:"id"`
	StartBlock        uint64       `json:"start_block"`
	EndBlock          uint64       `json:"end_block"`
	ValidatorSet      []*Validator `json:"validator_set"`
	SelectedProducers []string     `json:"selected_producers"`
	ChainId           string       `json:"chain_id"`
	CreatedAt         time.Time    `json:"created_at"`
}

// ProtoMessage 인터페이스 구현
func (m *Span) ProtoMessage() {}
func (m *Span) Reset()        { *m = Span{} }
func (m *Span) String() string {
	return fmt.Sprintf("Span{ID: %d, StartBlock: %d, EndBlock: %d, ChainID: %s}",
		m.Id, m.StartBlock, m.EndBlock, m.ChainId)
}

// NewSpan은 새로운 Span 객체를 생성합니다.
func NewSpan(
	id uint64,
	startBlock uint64,
	endBlock uint64,
	validatorSet []*Validator,
	selectedProducers []string,
	chainID string,
) *Span {
	return &Span{
		Id:                id,
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		ValidatorSet:      validatorSet,
		SelectedProducers: selectedProducers,
		ChainId:           chainID,
		CreatedAt:         time.Now(),
	}
}

// NewValidator는 새로운 Validator 객체를 생성합니다.
func NewValidator(
	address string,
	votingPower int64,
	proposerPriority int64,
) *Validator {
	return &Validator{
		Address:          address,
		VotingPower:      votingPower,
		ProposerPriority: proposerPriority,
	}
}
