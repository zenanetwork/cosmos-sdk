package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	checkpointkeeper "github.com/cosmos/cosmos-sdk/x/checkpoint/keeper"
	checkpointtypes "github.com/cosmos/cosmos-sdk/x/checkpoint/types"
	spankeeper "github.com/cosmos/cosmos-sdk/x/span/keeper"
	spantypes "github.com/cosmos/cosmos-sdk/x/span/types"
)

// IntegrationTestSuite는 checkpoint와 span 모듈 간의 통합 테스트를 위한 테스트 스위트입니다.
type IntegrationTestSuite struct {
	suite.Suite

	ctx                sdk.Context
	checkpointKeeper   checkpointkeeper.Keeper
	spanKeeper         spankeeper.Keeper
	cdc                codec.BinaryCodec
	checkpointStoreKey storetypes.StoreKey
	spanStoreKey       storetypes.StoreKey
}

// TestIntegrationTestSuite는 IntegrationTestSuite를 실행합니다.
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// SetupTest는 각 테스트 전에 실행되는 설정 함수입니다.
func (suite *IntegrationTestSuite) SetupTest() {
	// 테스트 환경 설정 코드
	// 실제 구현에서는 테스트 컨텍스트, 스토어, 키퍼 등을 설정해야 합니다.
}

// TestCheckpointWithSpan은 checkpoint 모듈이 span 모듈의 데이터를 사용하는 기능을 테스트합니다.
func (suite *IntegrationTestSuite) TestCheckpointWithSpan() {
	// 1. span 생성
	span := suite.spanKeeper.CreateSpan(
		suite.ctx,
		1,
		100,
		200,
		[]*spantypes.Validator{},
		[]string{},
		"test-chain",
	)

	// 현재 스팬 ID 설정
	suite.spanKeeper.SetCurrentSpanID(suite.ctx, span.Id)

	// 2. 체크포인트 생성
	checkpoint := suite.checkpointKeeper.CreateCheckpoint(
		suite.ctx,
		1,
		span.StartBlock,
		span.EndBlock,
		[]byte("test-root-hash"),
		"test-proposer",
	)

	// 3. 체크포인트 검증
	suite.Require().Equal(int64(1), checkpoint.Number, "체크포인트 번호가 일치해야 합니다")
	suite.Require().Equal(span.StartBlock, checkpoint.StartBlock, "시작 블록이 일치해야 합니다")
	suite.Require().Equal(span.EndBlock, checkpoint.EndBlock, "종료 블록이 일치해야 합니다")
	suite.Require().Equal("test-proposer", checkpoint.Proposer, "제안자가 일치해야 합니다")

	// 4. 현재 체크포인트 번호 설정
	suite.checkpointKeeper.SetCurrentCheckpointNumber(suite.ctx, checkpoint.Number)

	// 5. 최신 체크포인트 조회
	latestCheckpoint, err := suite.checkpointKeeper.GetLatestCheckpoint(suite.ctx)
	suite.Require().NoError(err, "최신 체크포인트를 조회할 때 오류가 없어야 합니다")
	suite.Require().Equal(checkpoint.Number, latestCheckpoint.Number, "체크포인트 번호가 일치해야 합니다")
}

// TestMultipleCheckpointsWithSpans은 여러 스팬과 체크포인트 간의 상호작용을 테스트합니다.
func (suite *IntegrationTestSuite) TestMultipleCheckpointsWithSpans() {
	// 1. 여러 스팬 생성
	spans := make([]*spantypes.Span, 3)

	spans[0] = suite.spanKeeper.CreateSpan(suite.ctx, 1, 100, 200, []*spantypes.Validator{}, []string{}, "test-chain")
	spans[1] = suite.spanKeeper.CreateSpan(suite.ctx, 2, 201, 300, []*spantypes.Validator{}, []string{}, "test-chain")
	spans[2] = suite.spanKeeper.CreateSpan(suite.ctx, 3, 301, 400, []*spantypes.Validator{}, []string{}, "test-chain")

	// 현재 스팬 ID 설정
	suite.spanKeeper.SetCurrentSpanID(suite.ctx, 3)

	// 2. 각 스팬에 대한 체크포인트 생성
	checkpoints := make([]*checkpointtypes.Checkpoint, 3)

	for i, span := range spans {
		checkpoints[i] = suite.checkpointKeeper.CreateCheckpoint(
			suite.ctx,
			int64(i+1),
			span.StartBlock,
			span.EndBlock,
			[]byte(fmt.Sprintf("hash-%d", i+1)),
			fmt.Sprintf("proposer-%d", i+1),
		)

		// 체크포인트 검증
		suite.Require().Equal(int64(i+1), checkpoints[i].Number, "체크포인트 번호가 일치해야 합니다")
		suite.Require().Equal(span.StartBlock, checkpoints[i].StartBlock, "시작 블록이 일치해야 합니다")
		suite.Require().Equal(span.EndBlock, checkpoints[i].EndBlock, "종료 블록이 일치해야 합니다")
	}

	// 3. 현재 체크포인트 번호 설정
	suite.checkpointKeeper.SetCurrentCheckpointNumber(suite.ctx, 3)

	// 4. 체크포인트 수 확인
	count := suite.checkpointKeeper.GetCheckpointCount(suite.ctx)
	suite.Require().Equal(int64(4), count, "체크포인트 수는 4여야 합니다") // 0번 체크포인트 + 3개 생성

	// 5. 최신 체크포인트 조회
	latestCheckpoint, err := suite.checkpointKeeper.GetLatestCheckpoint(suite.ctx)
	suite.Require().NoError(err, "최신 체크포인트를 조회할 때 오류가 없어야 합니다")
	suite.Require().Equal(int64(3), latestCheckpoint.Number, "최신 체크포인트 번호가 3이어야 합니다")
}
