package keeper

import (
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k *Keeper) CollectWager(ctx sdk.Context, storedGame *types.StoredGame) error {
	if storedGame.MoveCount == 0 {
		black, err := storedGame.GetBlackAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, black, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrBlackCannotPay.Error())
		}
	} else if storedGame.MoveCount == 1 {
		red, err := storedGame.GetRedAddress()
		if err != nil {
			panic(err.Error())
		}
		err = k.bank.SendCoinsFromAccountToModule(ctx, red, types.ModuleName, sdk.NewCoins(storedGame.GetWagerCoin()))
		if err != nil {
			return sdkerrors.Wrapf(err, types.ErrRedCannotPay.Error())
		}
	}
	return nil
}

func (k *Keeper) MustPayWinnings(ctx sdk.Context, storedGame *types.StoredGame) {

}

func (k *Keeper) MustRefundWager(ctx sdk.Context, storedGame *types.StoredGame) {

}
