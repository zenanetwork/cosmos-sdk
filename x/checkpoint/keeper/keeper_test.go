package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/keeper"
	"github.com/cosmos/cosmos-sdk/x/checkpoint/types"
)

// KeeperTestSuite는 checkpoint 모듈의 keeper를 테스트하기 위한 테스트 스위트입니다.
type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	keeper        keeper.Keeper
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	spanKeeper    types.SpanKeeper
}

// TestKeeperTestSuite는 KeeperTestSuite를 실행합니다.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest는 각 테스트 전에 실행되는 설정 함수입니다.
func (suite *KeeperTestSuite) SetupTest() {
	// 테스트 환경 설정 코드
	// 실제 구현에서는 테스트 컨텍스트, 스토어, 키퍼 등을 설정해야 합니다.
}

// TestGetSetCurrentCheckpointNumber는 GetCurrentCheckpointNumber와 SetCurrentCheckpointNumber 메서드를 테스트합니다.
func (suite *KeeperTestSuite) TestGetSetCurrentCheckpointNumber() {
	// 초기값 확인
	currentNumber := suite.keeper.GetCurrentCheckpointNumber(suite.ctx)
	suite.Require().Equal(int64(0), currentNumber, "초기 CurrentCheckpointNumber는 0이어야 합니다")

	// 값 설정 및 확인
	testNumber := int64(123)
	suite.keeper.SetCurrentCheckpointNumber(suite.ctx, testNumber)

	// 설정된 값 확인
	currentNumber = suite.keeper.GetCurrentCheckpointNumber(suite.ctx)
	suite.Require().Equal(testNumber, currentNumber, "설정된 CurrentCheckpointNumber가 일치해야 합니다")
}

// TestGetSetCheckpoint는 GetCheckpoint와 SetCheckpoint 메서드를 테스트합니다.
func (suite *KeeperTestSuite) TestGetSetCheckpoint() {
	// 테스트 체크포인트 생성
	testCheckpoint := types.NewCheckpoint(
		1,
		100,
		200,
		[]byte("test-root-hash"),
		"test-proposer",
	)

	// 체크포인트 저장
	suite.keeper.SetCheckpoint(suite.ctx, testCheckpoint)

	// 저장된 체크포인트 조회
	retrievedCheckpoint, err := suite.keeper.GetCheckpoint(suite.ctx, testCheckpoint.Number)
	suite.Require().NoError(err, "저장된 체크포인트를 조회할 때 오류가 없어야 합니다")
	suite.Require().Equal(testCheckpoint.Number, retrievedCheckpoint.Number, "체크포인트 번호가 일치해야 합니다")
	suite.Require().Equal(testCheckpoint.StartBlock, retrievedCheckpoint.StartBlock, "시작 블록이 일치해야 합니다")
	suite.Require().Equal(testCheckpoint.EndBlock, retrievedCheckpoint.EndBlock, "종료 블록이 일치해야 합니다")
	suite.Require().Equal(testCheckpoint.Proposer, retrievedCheckpoint.Proposer, "제안자가 일치해야 합니다")
}

// TestCreateCheckpoint는 CreateCheckpoint 메서드를 테스트합니다.
func (suite *KeeperTestSuite) TestCreateCheckpoint() {
	// 체크포인트 생성
	checkpoint := suite.keeper.CreateCheckpoint(
		suite.ctx,
		1,
		100,
		200,
		[]byte("test-root-hash"),
		"test-proposer",
	)

	// 생성된 체크포인트 확인
	suite.Require().Equal(int64(1), checkpoint.Number, "체크포인트 번호가 일치해야 합니다")
	suite.Require().Equal(uint64(100), checkpoint.StartBlock, "시작 블록이 일치해야 합니다")
	suite.Require().Equal(uint64(200), checkpoint.EndBlock, "종료 블록이 일치해야 합니다")
	suite.Require().Equal("test-proposer", checkpoint.Proposer, "제안자가 일치해야 합니다")

	// 저장된 체크포인트 조회
	retrievedCheckpoint, err := suite.keeper.GetCheckpoint(suite.ctx, checkpoint.Number)
	suite.Require().NoError(err, "생성된 체크포인트를 조회할 때 오류가 없어야 합니다")
	suite.Require().Equal(checkpoint.Number, retrievedCheckpoint.Number, "체크포인트 번호가 일치해야 합니다")
}

// TestGetLatestCheckpoint는 GetLatestCheckpoint 메서드를 테스트합니다.
func (suite *KeeperTestSuite) TestGetLatestCheckpoint() {
	// 여러 체크포인트 생성
	suite.keeper.CreateCheckpoint(suite.ctx, 1, 100, 200, []byte("hash1"), "proposer1")
	suite.keeper.CreateCheckpoint(suite.ctx, 2, 201, 300, []byte("hash2"), "proposer2")
	suite.keeper.CreateCheckpoint(suite.ctx, 3, 301, 400, []byte("hash3"), "proposer3")

	// 현재 체크포인트 번호 설정
	suite.keeper.SetCurrentCheckpointNumber(suite.ctx, 3)

	// 최신 체크포인트 조회
	latestCheckpoint, err := suite.keeper.GetLatestCheckpoint(suite.ctx)
	suite.Require().NoError(err, "최신 체크포인트를 조회할 때 오류가 없어야 합니다")
	suite.Require().Equal(int64(3), latestCheckpoint.Number, "최신 체크포인트 번호가 3이어야 합니다")
	suite.Require().Equal(uint64(301), latestCheckpoint.StartBlock, "시작 블록이 301이어야 합니다")
	suite.Require().Equal(uint64(400), latestCheckpoint.EndBlock, "종료 블록이 400이어야 합니다")
	suite.Require().Equal("proposer3", latestCheckpoint.Proposer, "제안자가 proposer3이어야 합니다")
}

// TestGetCheckpointCount는 GetCheckpointCount 메서드를 테스트합니다.
func (suite *KeeperTestSuite) TestGetCheckpointCount() {
	// 초기 체크포인트 수 확인
	count := suite.keeper.GetCheckpointCount(suite.ctx)
	suite.Require().Equal(int64(1), count, "초기 체크포인트 수는 1이어야 합니다")

	// 체크포인트 생성 및 현재 번호 설정
	suite.keeper.CreateCheckpoint(suite.ctx, 1, 100, 200, []byte("hash1"), "proposer1")
	suite.keeper.CreateCheckpoint(suite.ctx, 2, 201, 300, []byte("hash2"), "proposer2")
	suite.keeper.SetCurrentCheckpointNumber(suite.ctx, 2)

	// 체크포인트 수 확인
	count = suite.keeper.GetCheckpointCount(suite.ctx)
	suite.Require().Equal(int64(3), count, "체크포인트 수는 3이어야 합니다")
}
