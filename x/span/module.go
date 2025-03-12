package span

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"cosmossdk.io/core/registry"

	"github.com/zenanetwork/cosmos-sdk/client"
	"github.com/zenanetwork/cosmos-sdk/codec"
	sdk "github.com/zenanetwork/cosmos-sdk/types"
	"github.com/zenanetwork/cosmos-sdk/types/module"
	"github.com/zenanetwork/cosmos-sdk/x/span/client/cli"
	"github.com/zenanetwork/cosmos-sdk/x/span/keeper"
	"github.com/zenanetwork/cosmos-sdk/x/span/types"
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
func (AppModuleBasic) RegisterLegacyAminoCodec(registrar registry.AminoRegistrar) {
	types.RegisterLegacyAminoCodec(registrar)
}

// RegisterInterfaces는 모듈의 인터페이스를 인터페이스 레지스트리에 등록합니다.
func (AppModuleBasic) RegisterInterfaces(registrar registry.InterfaceRegistrar) {
	types.RegisterInterfaces(registrar)
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
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
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
}

// NewAppModule는 새로운 AppModule을 생성합니다.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

// Name은 모듈의 이름을 반환합니다.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices는 모듈의 서비스를 등록합니다.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants는 모듈의 불변성을 등록합니다.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis는 제네시스 상태를 초기화합니다.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	am.keeper.InitGenesis(ctx, genState)
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
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 현재 스팬 ID 확인
	currentSpanID := am.keeper.GetCurrentSpanID(sdkCtx)
	if currentSpanID == 0 {
		// 초기 스팬이 없는 경우 생성
		params := am.keeper.GetParams(sdkCtx)
		span := am.keeper.CreateSpan(
			sdkCtx,
			1,
			0,
			params.SpanLength,
			[]*types.Validator{},
			[]string{},
			params.ChainID,
		)
		am.keeper.SetCurrentSpanID(sdkCtx, span.Id)
	}

	return nil
}

// EndBlock은 블록 종료 시 호출됩니다.
func (am AppModule) EndBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// 현재 블록 높이 확인
	height := uint64(sdkCtx.BlockHeight())

	// 현재 스팬 확인
	currentSpanID := am.keeper.GetCurrentSpanID(sdkCtx)
	currentSpan, found := am.keeper.GetSpan(sdkCtx, currentSpanID)
	if !found {
		return fmt.Errorf("current span not found: %d", currentSpanID)
	}

	// 스팬 종료 확인
	if height >= currentSpan.EndBlock {
		// 새로운 스팬 생성
		params := am.keeper.GetParams(sdkCtx)
		newSpanID := currentSpanID + 1
		newSpan := am.keeper.CreateSpan(
			sdkCtx,
			newSpanID,
			currentSpan.EndBlock+1,
			currentSpan.EndBlock+params.SpanLength,
			[]*types.Validator{}, // 검증자 세트는 외부에서 업데이트
			[]string{},           // 생산자 목록은 외부에서 업데이트
			params.ChainID,
		)
		am.keeper.SetCurrentSpanID(sdkCtx, newSpan.Id)

		// 이벤트 발행
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeNewSpan,
				sdk.NewAttribute(types.AttributeKeySpanID, fmt.Sprintf("%d", newSpan.Id)),
				sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprintf("%d", newSpan.StartBlock)),
				sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprintf("%d", newSpan.EndBlock)),
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
