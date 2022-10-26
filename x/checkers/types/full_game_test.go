package types_test

import (
	"fmt"
	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const (
	alice      = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob        = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	badAddress = "THIS_IS_A_BAD_ADDRESS"
)

func GetStoredGame1() types.StoredGame {
	return types.StoredGame{
		Black:    alice,
		Red:      bob,
		Index:    "1",
		Turn:     "b",
		Board:    rules.New().String(),
		Deadline: types.DeadlineLayout,
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	expectedAddress, expectedAddrErr := sdk.AccAddressFromBech32(alice)
	actualAddress, actualAddrErr := GetStoredGame1().GetBlackAddress()
	require.Equal(t, expectedAddress, actualAddress)
	require.Nil(t, expectedAddrErr)
	require.Nil(t, actualAddrErr)
}

func TestGetAddressWrongBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = badAddress
	actualAddress, actualAddrErr := storedGame.GetBlackAddress()
	require.Nil(t, actualAddress)
	require.EqualError(
		t,
		actualAddrErr,
		fmt.Sprintf(
			"black address is invalid: %s: decoding bech32 failed: invalid separator index -1",
			badAddress),
	)
	require.EqualError(t, storedGame.Validate(), actualAddrErr.Error())
}

func TestCanGetAddressRed(t *testing.T) {
	expectedAddress, expectedAddrErr := sdk.AccAddressFromBech32(bob)
	actualAddress, actualAddrErr := GetStoredGame1().GetRedAddress()
	require.Equal(t, expectedAddress, actualAddress)
	require.Nil(t, expectedAddrErr)
	require.Nil(t, actualAddrErr)
}

func TestGetAddressWrongRed(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = badAddress
	actualAddress, actualAddrErr := storedGame.GetRedAddress()
	require.Nil(t, actualAddress)
	require.EqualError(
		t,
		actualAddrErr,
		fmt.Sprintf(
			"red address is invalid: %s: decoding bech32 failed: invalid separator index -1",
			badAddress),
	)
	require.EqualError(t, storedGame.Validate(), actualAddrErr.Error())
}

func TestParseGameCorrect(t *testing.T) {
	game, err := GetStoredGame1().ParseGame()
	require.EqualValues(t, rules.New(), game)
	require.Nil(t, err)
}

func TestParseGameIfChangedOk(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Board = strings.Replace(storedGame.Board, "b", "r", 1)
	game, err := storedGame.ParseGame()
	require.NotEqualValues(t, rules.New(), game)
	require.Nil(t, err)
}

func TestParseGameWrongPieceColor(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Board = strings.Replace(storedGame.Board, "b", "w", 1)
	game, err := storedGame.ParseGame()
	require.Nil(t, game)
	require.NotNil(t, err)
	require.EqualError(t, err, "game cannot be parsed: invalid board, invalid piece at 1, 0")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestParseGameWrongTurnColor(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Turn = "w"
	game, err := storedGame.ParseGame()
	require.Nil(t, game)
	require.NotNil(t, err)
	require.EqualError(t, err, fmt.Sprintf("game cannot be parsed: Turn: %s", storedGame.Turn))
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestParseDeadlineCorrect(t *testing.T) {
	deadline, err := GetStoredGame1().GetDeadlineAsTime()
	require.Nil(t, err)
	require.Equal(t, time.Time(time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.UTC)), deadline)
}

func TestParseDeadlineMissingMonth(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Deadline = "2006-02 15:04:05.999999999 +0000 UTC"
	_, err := storedGame.GetDeadlineAsTime()
	require.EqualError(
		t,
		err,
		"deadline cannot be parsed: 2006-02 15:04:05.999999999 +0000 UTC: parsing time \"2006-02 15:04:05.999999999 "+
			"+0000 UTC\" as \"2006-01-02 15:04:05.999999999 +0000 UTC\": cannot parse \" 15:04:05.999999999 +0000 "+
			"UTC\" as \"-\"")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestGetPlayerAddressBlackCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	black, found, err := storedGame.GetPlayerAddress("b")
	require.EqualValues(t, alice, black.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressBlackIncorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "notanaddress"
	black, found, err := storedGame.GetPlayerAddress("b")
	require.Nil(t, black)
	require.False(t, found)
	require.EqualError(
		t,
		err,
		"black address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1")
}

func TestGetPlayerAddressRedCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	red, found, err := storedGame.GetPlayerAddress("r")
	require.EqualValues(t, bob, red.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressRedIncorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "notanaddress"
	red, found, err := storedGame.GetPlayerAddress("r")
	require.Nil(t, red)
	require.False(t, found)
	require.EqualError(
		t,
		err,
		"red address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1")
}

func TestGetPlayerAddressWhiteNotFound(t *testing.T) {
	storedGame := GetStoredGame1()
	white, found, err := storedGame.GetPlayerAddress("w")
	require.Nil(t, white)
	require.False(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressNoPlayerNotFound(t *testing.T) {
	storedGame := GetStoredGame1()
	noPlayer, found, err := storedGame.GetPlayerAddress("*")
	require.Nil(t, noPlayer)
	require.False(t, found)
	require.Nil(t, err)
}

func TestGetWinnerAddressBlackCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Winner = "b"
	winner, found, err := storedGame.GetWinnerAddress()
	require.EqualValues(t, alice, winner.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetWinnerAddressRedCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Winner = "r"
	winner, found, err := storedGame.GetWinnerAddress()
	require.EqualValues(t, bob, winner.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetWinnerAddressNoWinnerYet(t *testing.T) {
	storedGame := GetStoredGame1()
	winner, found, err := storedGame.GetWinnerAddress()
	require.Nil(t, winner)
	require.False(t, found)
	require.Nil(t, err)
}

func TestValidateOk(t *testing.T) {
	storedGame := GetStoredGame1()
	require.NoError(t, storedGame.Validate())
}
