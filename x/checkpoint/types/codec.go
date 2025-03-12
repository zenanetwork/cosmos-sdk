package types

import (
	"context"

	"cosmossdk.io/core/registry"
	coretransaction "cosmossdk.io/core/transaction"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// MsgServer는 Msg 서비스를 위한 서버 API를 정의합니다.
type MsgServer interface {
	// CreateCheckpoint는 새로운 체크포인트를 생성합니다.
	CreateCheckpoint(context.Context, *MsgCreateCheckpoint) (*MsgCreateCheckpointResponse, error)
	// UpdateParams는 모듈 파라미터를 업데이트합니다.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// RegisterLegacyAminoCodec은 모듈의 인터페이스를 레거시 아미노 코덱에 등록합니다.
func RegisterLegacyAminoCodec(registrar registry.AminoRegistrar) {
	legacy.RegisterAminoMsg(registrar, &MsgCreateCheckpoint{}, "checkpoint/CreateCheckpoint")
	legacy.RegisterAminoMsg(registrar, &MsgUpdateParams{}, "checkpoint/UpdateParams")
}

// RegisterInterfaces는 모듈의 인터페이스를 인터페이스 레지스트리에 등록합니다.
func RegisterInterfaces(registrar registry.InterfaceRegistrar) {
	registrar.RegisterImplementations(
		(*coretransaction.Msg)(nil),
		&MsgCreateCheckpoint{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}

// _Msg_serviceDesc는 Msg 서비스 디스크립터를 정의합니다.
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "zenachain.checkpoint.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCheckpoint",
			Handler:    nil,
		},
		{
			MethodName: "UpdateParams",
			Handler:    nil,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "zenachain/checkpoint/v1/tx.proto",
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
}
