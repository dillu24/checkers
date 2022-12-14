package keeper_test

import (
	goContext "context"
	"fmt"
	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/testutil/mock_types"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func setupMsgServerWithOneGameForPlayMove(t testing.TB) (types.MsgServer, keeper.Keeper, goContext.Context,
	*gomock.Controller, *mock_types.MockBankEscrowKeeper) {
	ctrl := gomock.NewController(t)
	bankMock := mock_types.NewMockBankEscrowKeeper(ctrl)
	k, ctx := keepertest.CheckersKeeperWithMocks(t, bankMock)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	server := keeper.NewMsgServerImpl(*k)
	context := sdk.WrapSDKContext(ctx)
	server.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
		Wager:   45,
		Denom:   "stake",
	})
	return server, *k, context, ctrl, bankMock
}

func TestPlayMove(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgPlayMoveResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    rules.PieceStrings[rules.NO_PLAYER],
	}, *playMoveResponse)
}

func TestPlayMoveGameNotFound(t *testing.T) {
	msgServer, _, context, ctrl, _ := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	gameIndex := "2"
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: gameIndex,
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, playMoveResponse)
	require.Error(t, err, fmt.Sprintf("%s: game by id not found", gameIndex))
}

func TestPlayMoveSameBlackRed(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     bob,
		Wager:   46,
		Denom:   "stake",
	})
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "2",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgPlayMoveResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    rules.PieceStrings[rules.NO_PLAYER],
	}, *playMoveResponse)
}

func TestPlayMoveSavedGame(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{NextId: 2, FifoHeadIndex: "1", FifoTailIndex: "1"}, systemInfo)
	storedGame, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "r",
		Black:       bob,
		Red:         carol,
		MoveCount:   1,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "*",
		Wager:       45,
		Denom:       "stake",
	}, storedGame)
}

func TestPlayMoveEmitted(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	require.EqualValues(t, sdk.StringEvent{
		Type: types.MovePlayedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.MovePlayedEventCreator, Value: bob},
			{Key: types.MovePlayedEventGameIndex, Value: "1"},
			{Key: types.MovePlayedEventCapturedX, Value: "-1"},
			{Key: types.MovePlayedEventCapturedY, Value: "-1"},
			{Key: types.MovePlayedEventWinner, Value: rules.PieceStrings[rules.NO_PLAYER]},
			{Key: types.MovePlayedEventBoard,
				Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, events[0])
}

func TestPlayMoveCalledBank(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	escrow.ExpectPay(context, bob, 45).Times(1)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
}

func TestPlayMoveConsumedGas(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	before := ctx.GasMeter().GasConsumed()
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	after := ctx.GasMeter().GasConsumed()
	require.GreaterOrEqual(t, after, before+1_000)
}

func TestPlayMoveNotPlayer(t *testing.T) {
	msgServer, _, context, ctrl, _ := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   alice,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	require.Nil(t, playMoveResponse)
	require.Error(t, err, fmt.Sprintf("%s: message creator is not a player", alice))
}

func TestPlayMoveCannotParseGame(t *testing.T) {
	msgServer, k, context, ctrl, _ := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	storedGame, _ := k.GetStoredGame(ctx, "1")
	storedGame.Board = "invalid board"
	k.SetStoredGame(ctx, storedGame)
	defer func() {
		r := recover()
		require.NotNil(t, r, "The code did not panic")
		require.Equal(t, r, "game cannot be parsed: invalid board string: invalid board")
	}()
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
}

func TestPlayMoveWrongOutOfTurn(t *testing.T) {
	msgServer, _, context, ctrl, _ := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	require.Nil(t, playMoveResponse)
	require.Error(t, err, fmt.Sprintf("%s: player tried to play out of turn", carol))
}

func TestPlayMoveWrongPieceAtDestination(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     0,
		ToX:       0,
		ToY:       0,
	})
	require.Nil(t, playMoveResponse)
	require.Error(t, err, fmt.Sprintf("Invalid move: %v to %v: wrong move", rules.Pos{
		X: 1,
		Y: 0,
	}, rules.Pos{
		X: 0,
		Y: 0,
	}))
}

func TestPlayMove2(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	require.Nil(t, err)
	require.EqualValues(t, playMoveResponse, &types.MsgPlayMoveResponse{
		CapturedX: -1,
		CapturedY: -1,
		Winner:    rules.PieceStrings[rules.NO_PLAYER],
	})
}

func TestPlayMove2SavedGame(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
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
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{NextId: 2, FifoHeadIndex: "1", FifoTailIndex: "1"}, systemInfo)
	storedGame, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   2,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "*",
		Wager:       45,
		Denom:       "stake",
	}, storedGame)
}

func TestPlayMove2Emitted(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
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
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	require.EqualValues(t, []sdk.Attribute{
		{Key: types.MovePlayedEventCreator, Value: carol},
		{Key: types.GameCreatedEventGameIndex, Value: "1"},
		{Key: types.MovePlayedEventCapturedX, Value: "-1"},
		{Key: types.MovePlayedEventCapturedY, Value: "-1"},
		{Key: types.MovePlayedEventWinner, Value: rules.PieceStrings[rules.NO_PLAYER]},
		{Key: types.MovePlayedEventBoard,
			Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
	}, events[0].Attributes[6:])
}

func TestPlayMove2CalledBank(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	payBob := escrow.ExpectPay(context, bob, 45).Times(1)
	escrow.ExpectPay(context, carol, 45).Times(1).After(payBob)
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
}

func TestPlayMove3(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
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
	playMoveResponse, err := msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     2,
		FromY:     3,
		ToX:       0,
		ToY:       5,
	})
	require.Nil(t, err)
	require.EqualValues(t, playMoveResponse, &types.MsgPlayMoveResponse{
		CapturedX: 1,
		CapturedY: 4,
		Winner:    rules.PieceStrings[rules.NO_PLAYER],
	})
}

func TestPlayMove3SavedGame(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
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
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     2,
		FromY:     3,
		ToX:       0,
		ToY:       5,
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{NextId: 2, FifoHeadIndex: "1", FifoTailIndex: "1"}, systemInfo)
	storedGame, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "*b*b*b*b|b*b*b*b*|***b*b*b|********|********|b*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "r",
		Black:       bob,
		Red:         carol,
		MoveCount:   3,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "*",
		Wager:       45,
		Denom:       "stake",
	}, storedGame)
}

func TestPlayMove3CalledBank(t *testing.T) {
	msgServer, _, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	defer ctrl.Finish()
	payBob := escrow.ExpectPay(context, bob, 45).Times(1)
	escrow.ExpectPay(context, carol, 45).Times(1).After(payBob)
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
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     2,
		FromY:     3,
		ToX:       0,
		ToY:       5,
	})
}

func TestSavedPlayedDeadlineIsParseable(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	game, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	_, err := game.GetDeadlineAsTime()
	require.Nil(t, err)
}

func TestPlayerInfoNoAdditionOnNoWinner(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	aliceInfo, found := k.GetPlayerInfo(ctx, alice)
	require.False(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          "",
		WonCount:       0,
		LostCount:      0,
		ForfeitedCount: 0,
	}, aliceInfo)
	bobInfo, found := k.GetPlayerInfo(ctx, bob)
	require.False(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          "",
		WonCount:       0,
		LostCount:      0,
		ForfeitedCount: 0,
	}, bobInfo)
	carolInfo, found := k.GetPlayerInfo(ctx, carol)
	require.False(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          "",
		WonCount:       0,
		LostCount:      0,
		ForfeitedCount: 0,
	}, carolInfo)
}

func TestPlayerInfoNoUpdatedOnNoWinner(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	k.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          alice,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	})
	k.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	})
	k.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       7,
		LostCount:      8,
		ForfeitedCount: 9,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	aliceInfo, found := k.GetPlayerInfo(ctx, alice)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          alice,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	}, aliceInfo)
	bobInfo, found := k.GetPlayerInfo(ctx, bob)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          bob,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	}, bobInfo)
	carolInfo, found := k.GetPlayerInfo(ctx, carol)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          carol,
		WonCount:       7,
		LostCount:      8,
		ForfeitedCount: 9,
	}, carolInfo)
}

func TestLeaderboardNoAdditonOnNoWinner(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	leaderboard, found := k.GetLeaderboard(ctx)
	require.True(t, found)
	require.EqualValues(t, len(leaderboard.Winners), 0)
}

func TestLeaderboardNotUpdatedOnNoWinner(t *testing.T) {
	msgServer, k, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	escrow.ExpectAny(context)
	k.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          alice,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	})
	k.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	})
	k.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       7,
		LostCount:      8,
		ForfeitedCount: 9,
	})
	k.SetLeaderboard(ctx, types.Leaderboard{Winners: []types.WinningPlayer{
		{
			PlayerAddress: carol,
			WonCount:      7,
			DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
		},
		{
			PlayerAddress: bob,
			WonCount:      4,
			DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
		},
		{
			PlayerAddress: alice,
			WonCount:      1,
			DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
		},
	}})

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})

	leaderboard, found := k.GetLeaderboard(ctx)
	require.True(t, found)
	require.EqualValues(t, types.Leaderboard{Winners: []types.WinningPlayer{
		{
			PlayerAddress: carol,
			WonCount:      7,
			DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
		},
		{
			PlayerAddress: bob,
			WonCount:      4,
			DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
		},
		{
			PlayerAddress: alice,
			WonCount:      1,
			DateAdded:     "2006-01-02 15:05:06.999999999 +0000 UTC",
		},
	}}, leaderboard)
}
