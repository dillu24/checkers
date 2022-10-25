package keeper_test

import (
	goContext "context"
	"fmt"
	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupMsgServerWithOneGameForRejectGame(t testing.TB) (types.MsgServer, keeper.Keeper, goContext.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	return server, *k, context
}

func TestRejectGameWrongByCreator(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   alice,
		GameIndex: "1",
	})
	require.Nil(t, rejectGameResponse)
	require.Error(t, err, fmt.Sprintf("%s: message creator is not a player", alice))
}

func TestRejectGameByBlackNoMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgRejectGameResponse{}, *rejectGameResponse)
}

func TestRejectGameByBlackNoMoveRemovedGame(t *testing.T) {
	msgServer, k, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(
		t,
		types.SystemInfo{
			NextId:        2,
			FifoHeadIndex: types.NoFifoIndex,
			FifoTailIndex: types.NoFifoIndex,
		}, systemInfo)
	_, found = k.GetStoredGame(ctx, "1")
	require.False(t, found)
}

func TestRejectGameByBlackNoMoveEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	require.EqualValues(t, sdk.StringEvent{
		Type: types.GameRejectedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.GameRejectedEventCreator, Value: bob},
			{Key: types.GameRejectedEventGameIndex, Value: "1"},
		},
	}, events[0])
}

func TestRejectGameByRedNoMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgRejectGameResponse{}, *rejectGameResponse)
}

func TestRejectGameByRedNoMoveRemovedGame(t *testing.T) {
	msgServer, k, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(
		t,
		types.SystemInfo{
			NextId:        2,
			FifoHeadIndex: types.NoFifoIndex,
			FifoTailIndex: types.NoFifoIndex,
		}, systemInfo)
	_, found = k.GetStoredGame(ctx, "1")
	require.False(t, found)
}

func TestRejectGameByRedNoMoveEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	require.EqualValues(t, sdk.StringEvent{
		Type: types.GameRejectedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.GameRejectedEventCreator, Value: carol},
			{Key: types.GameRejectedEventGameIndex, Value: "1"},
		},
	}, events[0])
}

func TestRejectGameByRedOneMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgRejectGameResponse{}, *rejectGameResponse)
}

func TestRejectGameByRedOneMoveRemovedGame(t *testing.T) {
	msgServer, k, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(
		t,
		types.SystemInfo{
			NextId:        2,
			FifoHeadIndex: types.NoFifoIndex,
			FifoTailIndex: types.NoFifoIndex,
		}, systemInfo)
	_, found = k.GetStoredGame(ctx, "1")
	require.False(t, found)
}

func TestRejectGameByRedOneMoveEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	require.EqualValues(t, sdk.StringEvent{
		Type: types.GameRejectedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.GameRejectedEventCreator, Value: carol},
			{Key: types.GameRejectedEventGameIndex, Value: "1"},
		},
	}, events[0])
}

func TestRejectGameByBlackWrongOneMove(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   bob,
		GameIndex: "1",
	})
	require.Nil(t, rejectGameResponse)
	require.Error(t, err, types.ErrBlackAlreadyPlayed)
}

func TestRejectGameByRedWrongTwoMoves(t *testing.T) {
	msgServer, _, context := setupMsgServerWithOneGameForRejectGame(t)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	rejectGameResponse, err := msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	require.Nil(t, rejectGameResponse)
	require.Error(t, err, types.ErrRedAlreadyPlayed)
}
