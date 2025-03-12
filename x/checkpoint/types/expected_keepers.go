package types

import (
	sdk "github.com/zenanetwork/cosmos-sdk/types"
	"github.com/zenanetwork/cosmos-sdk/x/auth/types"
	spantypes "github.com/zenanetwork/cosmos-sdk/x/span/types"
)

// AccountKeeper는 계정 모듈의 인터페이스를 정의합니다.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// BankKeeper는 뱅크 모듈의 인터페이스를 정의합니다.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// SpanKeeper는 span 모듈의 인터페이스를 정의합니다.
type SpanKeeper interface {
	GetSpan(ctx sdk.Context, spanID uint64) (*spantypes.Span, bool)
	GetCurrentSpanID(ctx sdk.Context) uint64
	GetSpanByHeight(ctx sdk.Context, height uint64) (*spantypes.Span, bool)
}
