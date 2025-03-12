package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validator는 검증자 정보를 나타냅니다.
type Validator struct {
	Address          string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	VotingPower      int64  `protobuf:"varint,2,opt,name=voting_power,json=votingPower,proto3" json:"voting_power,omitempty"`
	ProposerPriority int64  `protobuf:"varint,3,opt,name=proposer_priority,json=proposerPriority,proto3" json:"proposer_priority,omitempty"`
}

// Span은 블록 범위와 관련 정보를 나타냅니다.
type Span struct {
	Id                uint64       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	StartBlock        uint64       `protobuf:"varint,2,opt,name=start_block,json=startBlock,proto3" json:"start_block,omitempty"`
	EndBlock          uint64       `protobuf:"varint,3,opt,name=end_block,json=endBlock,proto3" json:"end_block,omitempty"`
	ValidatorSet      []*Validator `protobuf:"bytes,4,rep,name=validator_set,json=validatorSet,proto3" json:"validator_set,omitempty"`
	SelectedProducers []string     `protobuf:"bytes,5,rep,name=selected_producers,json=selectedProducers,proto3" json:"selected_producers,omitempty"`
	ChainId           string       `protobuf:"bytes,6,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	CreatedAt         time.Time    `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
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
func NewValidator(address string, votingPower int64, proposerPriority int64) *Validator {
	return &Validator{
		Address:          address,
		VotingPower:      votingPower,
		ProposerPriority: proposerPriority,
	}
}

// ValidatorFromSDK는 SDK 검증자에서 Validator 객체를 생성합니다.
func ValidatorFromSDK(validator sdk.ValidatorI) *Validator {
	return &Validator{
		Address:          validator.GetOperator().String(),
		VotingPower:      validator.GetConsensusPower(sdk.DefaultPowerReduction),
		ProposerPriority: validator.GetProposerPriority(),
	}
}
