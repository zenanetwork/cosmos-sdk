package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateCheckpoint{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// ProtoMessage 인터페이스 구현
func (m *MsgCreateCheckpoint) ProtoMessage() {}
func (m *MsgCreateCheckpoint) Reset()        { *m = MsgCreateCheckpoint{} }
func (m *MsgCreateCheckpoint) String() string {
	return fmt.Sprintf("MsgCreateCheckpoint{Creator: %s, StartBlock: %d, EndBlock: %d}",
		m.Creator, m.StartBlock, m.EndBlock)
}

func (m *MsgUpdateParams) ProtoMessage() {}
func (m *MsgUpdateParams) Reset()        { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string {
	return fmt.Sprintf("MsgUpdateParams{Authority: %s, Params: %s}",
		m.Authority, m.Params.String())
}

// MsgCreateCheckpoint는 새로운 체크포인트를 생성하기 위한 메시지입니다.
type MsgCreateCheckpoint struct {
	Creator    string `json:"creator"`
	StartBlock uint64 `json:"start_block"`
	EndBlock   uint64 `json:"end_block"`
	RootHash   []byte `json:"root_hash"`
}

// NewMsgCreateCheckpoint은 새로운 MsgCreateCheckpoint 객체를 생성합니다.
func NewMsgCreateCheckpoint(
	creator string,
	startBlock uint64,
	endBlock uint64,
	rootHash []byte,
) *MsgCreateCheckpoint {
	return &MsgCreateCheckpoint{
		Creator:    creator,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   rootHash,
	}
}

// Route는 메시지 라우팅 이름을 반환합니다.
func (msg MsgCreateCheckpoint) Route() string {
	return RouterKey
}

// Type은 메시지 타입을 반환합니다.
func (msg MsgCreateCheckpoint) Type() string {
	return "CreateCheckpoint"
}

// ValidateBasic는 메시지의 기본 유효성을 검사합니다.
func (msg MsgCreateCheckpoint) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}

	if msg.StartBlock >= msg.EndBlock {
		return ErrInvalidBlockRange
	}

	if len(msg.RootHash) == 0 {
		return ErrInvalidRootHash
	}

	return nil
}

// GetSigners는 메시지 서명자를 반환합니다.
func (msg MsgCreateCheckpoint) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes는 메시지 서명 바이트를 반환합니다.
func (msg MsgCreateCheckpoint) GetSignBytes() []byte {
	bz, _ := ModuleCdc.MarshalJSON(&msg)
	return bz
}

// MsgCreateCheckpointResponse는 MsgCreateCheckpoint에 대한 응답입니다.
type MsgCreateCheckpointResponse struct {
	Number int64 `json:"number"`
}

// ProtoMessage 인터페이스 구현
func (m *MsgCreateCheckpointResponse) ProtoMessage() {}
func (m *MsgCreateCheckpointResponse) Reset()        { *m = MsgCreateCheckpointResponse{} }
func (m *MsgCreateCheckpointResponse) String() string {
	return fmt.Sprintf("MsgCreateCheckpointResponse{Number: %d}", m.Number)
}

// MsgUpdateParams는 모듈 파라미터를 업데이트하기 위한 메시지입니다.
type MsgUpdateParams struct {
	Authority string `json:"authority"`
	Params    Params `json:"params"`
}

// NewMsgUpdateParams은 새로운 MsgUpdateParams 객체를 생성합니다.
func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// Route는 메시지 라우팅 이름을 반환합니다.
func (msg MsgUpdateParams) Route() string {
	return RouterKey
}

// Type은 메시지 타입을 반환합니다.
func (msg MsgUpdateParams) Type() string {
	return "UpdateParams"
}

// ValidateBasic는 메시지의 기본 유효성을 검사합니다.
func (msg MsgUpdateParams) ValidateBasic() error {
	if msg.Authority == "" {
		return fmt.Errorf("authority cannot be empty")
	}

	return msg.Params.Validate()
}

// GetSigners는 메시지 서명자를 반환합니다.
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

// GetSignBytes는 메시지 서명 바이트를 반환합니다.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	bz, _ := ModuleCdc.MarshalJSON(&msg)
	return bz
}

// MsgUpdateParamsResponse는 MsgUpdateParams에 대한 응답입니다.
type MsgUpdateParamsResponse struct{}

// ProtoMessage 인터페이스 구현
func (m *MsgUpdateParamsResponse) ProtoMessage() {}
func (m *MsgUpdateParamsResponse) Reset()        { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string {
	return "MsgUpdateParamsResponse{}"
}
