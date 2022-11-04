package keeper_test

import (
	goContext "context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"

	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/testutil/mock_types"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
)

func setupKeeperForWagerHandler(t testing.TB) (keeper.Keeper, goContext.Context, *gomock.Controller,
	*mock_types.MockBankEscrowKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := mock_types.NewMockBankEscrowKeeper(ctrl)
	k, ctx := keepertest.CheckersKeeperWithMocks(t, bankMock)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	context := sdk.WrapSDKContext(ctx)
	return *k, context, ctrl, bankMock
}

func TestWagerHandlerCollectWrongNoBlack(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, "black address is invalid: : empty address string is not allowed", r)
	}()
	k.CollectWager(ctx, &types.StoredGame{MoveCount: 0})
}

func TestWagerHandlerCollectFailedNoMove(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	black, _ := sdk.AccAddressFromBech32(alice)
	escrow.EXPECT().
		SendCoinsFromAccountToModule(ctx, black, types.ModuleName, gomock.Any()).
		Return(errors.New("oops"))
	err := k.CollectWager(ctx, &types.StoredGame{
		Black:     alice,
		MoveCount: 0,
		Wager:     45,
		Denom:     "stake",
	})
	require.NotNil(t, err)
	require.EqualError(t, err, "black cannot pay the wager: oops")
}

func TestWagerHandlerCollectWrongNoRed(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, "red address is invalid: : empty address string is not allowed", r)
	}()
	k.CollectWager(ctx, &types.StoredGame{MoveCount: 1})
}

func TestWagerHandlerCollectFailedOneMove(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	red, _ := sdk.AccAddressFromBech32(bob)
	escrow.EXPECT().
		SendCoinsFromAccountToModule(ctx, red, types.ModuleName, gomock.Any()).
		Return(errors.New("oops"))
	err := k.CollectWager(ctx, &types.StoredGame{
		Red:       bob,
		MoveCount: 1,
		Wager:     45,
		Denom:     "stake",
	})
	require.NotNil(t, err)
	require.EqualError(t, err, "red cannot pay the wager: oops")
}

func TestWagerHandlerCollectNoMove(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectPay(context, alice, 45).Times(1)
	err := k.CollectWager(ctx, &types.StoredGame{
		Black:     alice,
		MoveCount: 0,
		Wager:     45,
		Denom:     "stake",
	})
	require.Nil(t, err)
}

func TestWagerHandlerCollectOneMove(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectPay(context, bob, 45).Times(1)
	err := k.CollectWager(ctx, &types.StoredGame{
		Red:       bob,
		MoveCount: 1,
		Wager:     45,
		Denom:     "stake",
	})
	require.Nil(t, err)
}

func TestWagerHandlerPayWrongNoWinnerAddress(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic!")
		require.Equal(t, "black address is invalid: : empty address string is not allowed", r)
	}()
	k.MustPayWinnings(ctx, &types.StoredGame{
		Winner: "b",
	})
}

func TestWagerHandlerPayWrongWinnerNotFound(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic!")
		require.Equal(t, "cannot find winner by color: *", r)
	}()
	k.MustPayWinnings(ctx, &types.StoredGame{
		Winner: "*",
		Red:    bob,
		Black:  alice,
	})
}

func TestWagerHandlerPayWrongNotPayTime(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, "there is nothing to pay, should not have been called", r)
	}()
	k.MustPayWinnings(ctx, &types.StoredGame{
		Black:     alice,
		Red:       bob,
		MoveCount: 0,
		Winner:    "b",
		Denom:     "stake",
	})
}

func TestWagerHandlerPayWrongEscrowFailed(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	black, _ := sdk.AccAddressFromBech32(alice)
	defer ctrl.Finish()
	escrow.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, black,
		gomock.Any()).Times(1).Return(errors.New("oops"))
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, "cannot pay winnings to winner: oops", r)
	}()
	k.MustPayWinnings(ctx, &types.StoredGame{
		Black:     alice,
		Red:       bob,
		MoveCount: 2,
		Winner:    "b",
		Wager:     45,
		Denom:     "stake",
	})
}

func TestWagerHandlerPayEscrowCalledOneMove(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectRefund(context, alice, 45).Times(1)
	k.MustPayWinnings(ctx, &types.StoredGame{
		Black:     alice,
		Red:       bob,
		MoveCount: 1,
		Winner:    "b",
		Wager:     45,
		Denom:     "stake",
	})
}

func TestWagerHandlerPayEscrowCalledTwoMoves(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectRefund(context, alice, 90).Times(1)
	k.MustPayWinnings(ctx, &types.StoredGame{
		Black:     alice,
		Red:       bob,
		MoveCount: 2,
		Winner:    "b",
		Wager:     45,
		Denom:     "stake",
	})
}

func TestWagerHandlerRefundWrongManyMoves(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic!")
		require.Equal(t, "game is not in a state to refund, move count: 2", r)
	}()
	k.MustRefundWager(ctx, &types.StoredGame{
		MoveCount: 2,
	})
}

func TestWagerHandlerRefundNoMoves(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	k.MustRefundWager(ctx, &types.StoredGame{MoveCount: 0})
}

func TestWagerHandlerRefundWrongNoBlack(t *testing.T) {
	k, context, ctrl, _ := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic!")
		require.EqualValues(t, "black address is invalid: : empty address string is not allowed", r)
	}()
	k.MustRefundWager(ctx, &types.StoredGame{MoveCount: 1})
}

func TestWagerHandlerRefundWrongEscrowFailed(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	black, _ := sdk.AccAddressFromBech32(alice)
	escrow.EXPECT().SendCoinsFromModuleToAccount(ctx, types.ModuleName, black,
		gomock.Any()).Times(1).Return(errors.New("oops"))
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic!")
		require.EqualValues(t, "cannot refund wager to: oops", r)
	}()
	k.MustRefundWager(ctx, &types.StoredGame{
		Black:     alice,
		MoveCount: 1,
		Wager:     45,
		Denom:     "stake",
	})
}

func TestWagerHandlerRefundCalled(t *testing.T) {
	k, context, ctrl, escrow := setupKeeperForWagerHandler(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectRefund(context, alice, 45).Times(1)
	k.MustRefundWager(ctx, &types.StoredGame{
		Black:     alice,
		MoveCount: 1,
		Wager:     45,
		Denom:     "stake",
	})
}
