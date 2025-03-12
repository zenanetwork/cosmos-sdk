package types

import (
	"context"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// MsgServer는 Msg 서비스를 위한 서버 API를 정의합니다.
type MsgServer interface {
	// CreateSpan은 새로운 스팬을 생성합니다.
	CreateSpan(context.Context, *MsgCreateSpan) (*MsgCreateSpanResponse, error)
	// UpdateParams는 모듈 파라미터를 업데이트합니다.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// RegisterLegacyAminoCodec은 모듈의 인터페이스를 레거시 아미노 코덱에 등록합니다.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateSpan{}, "span/CreateSpan")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "span/UpdateParams")
}

// RegisterInterfaces는 모듈의 인터페이스를 인터페이스 레지스트리에 등록합니다.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateSpan{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// _Msg_serviceDesc는 Msg 서비스 디스크립터를 정의합니다.
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cosmos.span.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSpan",
			Handler:    nil,
		},
		{
			MethodName: "UpdateParams",
			Handler:    nil,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cosmos/span/v1/tx.proto",
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
}
