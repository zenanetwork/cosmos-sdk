package checkpoint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/client/cli"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/keeper"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

// AppModuleBasic는 모듈의 기본 정보를 정의합니다.
type AppModuleBasic struct {
	cdc codec.Codec
}

// Name은 모듈의 이름을 반환합니다.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec은 모듈의 인터페이스를 레거시 아미노 코덱에 등록합니다.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces는 모듈의 인터페이스를 인터페이스 레지스트리에 등록합니다.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis는 기본 제네시스 상태를 반환합니다.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis는 제네시스 상태의 유효성을 검사합니다.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes는 gRPC 게이트웨이 라우트를 등록합니다.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// 필요한 경우 여기에 구현
}

// GetTxCmd는 모듈의 트랜잭션 명령어를 반환합니다.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd는 모듈의 쿼리 명령어를 반환합니다.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule는 모듈의 구현을 정의합니다.
type AppModule struct {
	AppModuleBasic

	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	spanKeeper    types.SpanKeeper
}

// NewAppModule는 새로운 AppModule을 생성합니다.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	spanKeeper types.SpanKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		spanKeeper:     spanKeeper,
	}
}

// Name은 모듈의 이름을 반환합니다.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices는 모듈의 서비스를 등록합니다.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// 필요한 경우 여기에 구현
}

// RegisterInvariants는 모듈의 불변성을 등록합니다.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis는 제네시스 상태를 초기화합니다.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	am.keeper.InitGenesis(ctx, &genState)
}

// ExportGenesis는 현재 상태를 제네시스 상태로 내보냅니다.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion은 모듈의 컨센서스 버전을 반환합니다.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock은 블록 시작 시 호출됩니다.
func (am AppModule) BeginBlock(ctx context.Context) error {
	// 체크포인트 관련 초기화 작업이 필요한 경우 여기에 구현
	return nil
}

// EndBlock은 블록 종료 시 호출됩니다.
func (am AppModule) EndBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 현재 블록 높이 확인
	height := uint64(sdkCtx.BlockHeight())

	// 체크포인트 생성 주기 확인
	params := am.keeper.GetParams(sdkCtx)
	if height > 0 && height%params.CheckpointBufferSize == 0 {
		// 현재 체크포인트 번호 가져오기
		currentNumber := am.keeper.GetCurrentCheckpointNumber(sdkCtx)
		newNumber := currentNumber + 1

		// 스팬 정보 확인 (간단한 검증만 수행)
		currentSpanID := am.spanKeeper.GetCurrentSpanID(sdkCtx)
		if currentSpanID == 0 {
			return fmt.Errorf("no active span found")
		}

		// 체크포인트 생성
		startBlock := height - params.CheckpointBufferSize + 1
		endBlock := height

		// 루트 해시는 실제 구현에서 계산해야 함
		// 여기서는 임시로 빈 바이트 배열 사용
		rootHash := []byte{}

		// 제안자는 현재 컨텍스트에서 가져와야 함
		// 여기서는 임시로 빈 문자열 사용
		proposer := ""

		checkpoint := am.keeper.CreateCheckpoint(
			sdkCtx,
			newNumber,
			startBlock,
			endBlock,
			rootHash,
			proposer,
		)

		// 현재 체크포인트 번호 업데이트
		am.keeper.SetCurrentCheckpointNumber(sdkCtx, newNumber)

		// 이벤트 발행
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeNewCheckpoint,
				sdk.NewAttribute(types.AttributeKeyCheckpointNumber, fmt.Sprintf("%d", checkpoint.Number)),
				sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprintf("%d", checkpoint.StartBlock)),
				sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprintf("%d", checkpoint.EndBlock)),
				sdk.NewAttribute(types.AttributeKeyProposer, checkpoint.Proposer),
			),
		)
	}

	return nil
}

// PrecommitFilter는 블록 커밋 전에 호출됩니다.
func (am AppModule) PrecommitFilter(ctx context.Context) error {
	// 현재 구현에서는 특별한 작업이 필요 없음
	return nil
}

// PrepareCheckState는 체크 상태 준비 시 호출됩니다.
func (am AppModule) PrepareCheckState(ctx context.Context) error {
	// 현재 구현에서는 특별한 작업이 필요 없음
	return nil
}

// IsAppModule은 AppModule 인터페이스를 구현합니다.
func (am AppModule) IsAppModule() {}

// IsOnePerModuleType은 AppModule 인터페이스를 구현합니다.
func (am AppModule) IsOnePerModuleType() {}
