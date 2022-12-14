package keeper_test

import (
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRejectFirstGameHasSavedFifo(t *testing.T) {
	msgServer, k, context, ctrl, _ := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   46,
		Denom:   "coin",
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "1",
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        3,
		FifoHeadIndex: "2",
		FifoTailIndex: "2",
	}, systemInfo)
	game2, found := k.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       carol,
		Red:         alice,
		MoveCount:   0,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "*",
		Wager:       46,
		Denom:       "coin",
	}, game2)
}

func TestRejectMiddleGameHasSaved(t *testing.T) {
	msgServer, k, context, ctrl, _ := setupMsgServerWithOneGameForRejectGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   46,
		Denom:   "coin",
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
		Wager:   47,
		Denom:   "gold",
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator:   carol,
		GameIndex: "2",
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        4,
		FifoHeadIndex: "1",
		FifoTailIndex: "3",
	}, systemInfo)
	game1, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   0,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  "3",
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "*",
		Wager:       45,
		Denom:       "stake",
	}, game1)
	game3, found := k.GetStoredGame(ctx, "3")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "3",
		Board:       "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:        "b",
		Black:       alice,
		Red:         bob,
		MoveCount:   0,
		BeforeIndex: "1",
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(ctx.BlockTime().Add(types.MaxTurnDuration)),
		Winner:      "*",
		Wager:       47,
		Denom:       "gold",
	}, game3)
}
