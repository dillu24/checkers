package types_test

import (
	"fmt"
	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	alice       = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob         = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	carol       = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	bad_address = "THIS_IS_A_BAD_ADDRESS"
)

func GetStoredGame1() types.StoredGame {
	return types.StoredGame{
		Black: alice,
		Red:   bob,
		Index: "1",
		Turn:  "b",
		Board: rules.New().String(),
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
	storedGame.Black = bad_address
	actualAddress, actualAddrErr := storedGame.GetBlackAddress()
	require.Nil(t, actualAddress)
	require.EqualError(
		t,
		actualAddrErr,
		fmt.Sprintf(
			"black address is invalid: %s: decoding bech32 failed: invalid separator index -1",
			bad_address),
	)
	require.EqualError(t, storedGame.Validate(), actualAddrErr.Error())
}
