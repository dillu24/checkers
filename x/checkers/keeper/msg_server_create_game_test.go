package keeper_test

import (
	goContext "context"
	"fmt"
	"testing"

	keepertest "github.com/alice/checkers/testutil/keeper"
	"github.com/alice/checkers/x/checkers"
	"github.com/alice/checkers/x/checkers/keeper"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	alice      = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob        = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	carol      = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	badAddress = "notAnAddress"
)

func setupMsgServerCreateGame(t testing.TB) (types.MsgServer, keeper.Keeper, goContext.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), *k, sdk.WrapSDKContext(ctx)
}

func TestCreateGame(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		GameIndex: "1",
	}, *createResponse)
}

func TestCreate1GameHasSaved(t *testing.T) {
	msgSrvr, k, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId: 2,
	}, systemInfo)
	game, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:     "1",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     bob,
		Red:       carol,
		MoveCount: 0,
	}, game)
}

func TestCreate1GameGetAll(t *testing.T) {
	msgSrvr, k, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	games := k.GetAllStoredGame(ctx)
	require.Len(t, games, 1)
	require.Equal(t, types.StoredGame{
		Index:     "1",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     bob,
		Red:       carol,
		MoveCount: 0,
	}, games[0])
}

func TestCreate1GameEmitted(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	require.EqualValues(t, sdk.StringEvent{
		Type: types.GameCreatedEventType,
		Attributes: []sdk.Attribute{
			{Key: types.GameCreatedEventCreator, Value: alice},
			{Key: types.GameCreatedEventGameIndex, Value: "1"},
			{Key: types.GameCreatedEventBlack, Value: bob},
			{Key: types.GameCreatedEventRed, Value: carol},
		},
	}, events[0])
}

func TestCreateGameRedAddressBad(t *testing.T) {
	msgSrvr, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     badAddress,
	})
	require.Nil(t, createResponse)
	require.Error(t,
		err,
		fmt.Sprintf("red address is invalid: %s: decoding bech32 failed: invalid separator index -1", badAddress),
	)
}

func TestCreateGameEmptyRedAddress(t *testing.T) {
	msgSrvr, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     "",
	})
	require.Nil(t, createResponse)
	require.Error(t,
		err,
		"red address is invalid: : empty address string is not allowed",
	)
}

func TestCreate3Games(t *testing.T) {
	msgSrvr, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "1",
	}, createResponse)
	createResponse, err = msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "2",
	}, createResponse)
	createResponse, err = msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})
	require.Nil(t, err)
	require.EqualValues(t, &types.MsgCreateGameResponse{
		GameIndex: "3",
	}, createResponse)
}

func TestCreate3GamesHasSaved(t *testing.T) {
	msgSrvr, k, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{NextId: 2}, systemInfo)
	storedGame, found := k.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:     "1",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     bob,
		Red:       carol,
		MoveCount: 0,
	}, storedGame)

	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
	})
	systemInfo, found = k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{NextId: 3}, systemInfo)
	storedGame, found = k.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:     "2",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     alice,
		Red:       bob,
		MoveCount: 0,
	}, storedGame)

	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})
	systemInfo, found = k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{NextId: 4}, systemInfo)
	storedGame, found = k.GetStoredGame(ctx, "3")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:     "3",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     carol,
		Red:       alice,
		MoveCount: 0,
	}, storedGame)
}

func TestCreate3GamesGetAll(t *testing.T) {
	msgSrvr, k, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
	})
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
	})
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
	})
	games := k.GetAllStoredGame(ctx)
	require.Len(t, games, 3)
	require.EqualValues(t, types.StoredGame{
		Index:     "1",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     bob,
		Red:       carol,
		MoveCount: 0,
	}, games[0])
	require.EqualValues(t, types.StoredGame{
		Index:     "2",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     alice,
		Red:       bob,
		MoveCount: 0,
	}, games[1])
	require.EqualValues(t, types.StoredGame{
		Index:     "3",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     carol,
		Red:       alice,
		MoveCount: 0,
	}, games[2])
}

func TestCreateGameFarFuture(t *testing.T) {
	msgSrvr, k, context := setupMsgServerCreateGame(t)
	ctx := sdk.UnwrapSDKContext(context)
	systemInfo, _ := k.GetSystemInfo(ctx)
	systemInfo.NextId = 1024
	k.SetSystemInfo(ctx, systemInfo)
	createResponse, err := msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   bob,
		Red:     alice,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{GameIndex: "1024"}, *createResponse)
	systemInfo, found := k.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId: 1025,
	}, systemInfo)
	storedGame, found := k.GetStoredGame(ctx, "1024")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:     "1024",
		Board:     "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Black:     bob,
		Red:       alice,
		MoveCount: 0,
	}, storedGame)
}
