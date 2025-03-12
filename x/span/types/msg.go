package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateSpan{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// ProtoMessage 인터페이스 구현
func (m *MsgCreateSpan) ProtoMessage() {}
func (m *MsgCreateSpan) Reset()        { *m = MsgCreateSpan{} }
func (m *MsgCreateSpan) String() string {
	return fmt.Sprintf("MsgCreateSpan{Creator: %s, StartBlock: %d, EndBlock: %d}",
		m.Creator, m.StartBlock, m.EndBlock)
}

func (m *MsgUpdateParams) ProtoMessage() {}
func (m *MsgUpdateParams) Reset()        { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string {
	return fmt.Sprintf("MsgUpdateParams{Authority: %s, Params: %s}",
		m.Authority, m.Params.String())
}

// MsgCreateSpan은 새로운 스팬을 생성하기 위한 메시지입니다.
type MsgCreateSpan struct {
	Creator           string       `json:"creator"`
	StartBlock        uint64       `json:"start_block"`
	EndBlock          uint64       `json:"end_block"`
	Validators        []*Validator `json:"validators"`
	SelectedProducers []string     `json:"selected_producers"`
	ChainId           string       `json:"chain_id"`
}

// NewMsgCreateSpan은 새로운 MsgCreateSpan 객체를 생성합니다.
func NewMsgCreateSpan(
	creator string,
	startBlock uint64,
	endBlock uint64,
	validators []*Validator,
	selectedProducers []string,
	chainID string,
) *MsgCreateSpan {
	return &MsgCreateSpan{
		Creator:           creator,
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		Validators:        validators,
		SelectedProducers: selectedProducers,
		ChainId:           chainID,
	}
}

// Route는 메시지 라우팅 이름을 반환합니다.
func (msg MsgCreateSpan) Route() string {
	return RouterKey
}

// Type은 메시지 타입을 반환합니다.
func (msg MsgCreateSpan) Type() string {
	return "CreateSpan"
}

// GetSigners는 메시지 서명자를 반환합니다.
func (msg MsgCreateSpan) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes는 메시지 서명 바이트를 반환합니다.
func (msg MsgCreateSpan) GetSignBytes() []byte {
	bz, _ := ModuleCdc.MarshalJSON(&msg)
	return bz
}

// ValidateBasic는 메시지의 기본 유효성을 검사합니다.
func (msg MsgCreateSpan) ValidateBasic() error {
	if msg.Creator == "" {
		return fmt.Errorf("creator cannot be empty")
	}

	if msg.StartBlock >= msg.EndBlock {
		return ErrInvalidBlockRange
	}

	if len(msg.Validators) == 0 {
		return ErrInvalidValidatorSet
	}

	return nil
}

// MsgCreateSpanResponse는 MsgCreateSpan에 대한 응답입니다.
type MsgCreateSpanResponse struct {
	Id uint64 `json:"id"`
}

// ProtoMessage 인터페이스 구현
func (m *MsgCreateSpanResponse) ProtoMessage() {}
func (m *MsgCreateSpanResponse) Reset()        { *m = MsgCreateSpanResponse{} }
func (m *MsgCreateSpanResponse) String() string {
	return fmt.Sprintf("MsgCreateSpanResponse{Id: %d}", m.Id)
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

// ValidateBasic는 메시지의 기본 유효성을 검사합니다.
func (msg MsgUpdateParams) ValidateBasic() error {
	if msg.Authority == "" {
		return fmt.Errorf("authority cannot be empty")
	}

	return msg.Params.Validate()
}

// MsgUpdateParamsResponse는 MsgUpdateParams에 대한 응답입니다.
type MsgUpdateParamsResponse struct{}

// ProtoMessage 인터페이스 구현
func (m *MsgUpdateParamsResponse) ProtoMessage()  {}
func (m *MsgUpdateParamsResponse) Reset()         { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string { return "MsgUpdateParamsResponse{}" }
