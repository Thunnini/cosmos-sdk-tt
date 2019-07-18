package stakingmint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingibc "github.com/cosmos/cosmos-sdk/x/staking-ibc"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type StakingMintKeeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	ibcKeeper     mock.Keeper
	stakingKeeper staking.Keeper
	supplyKeeper  supply.Keeper
}

func NewStakingMintKeeper(cdc *codec.Codec, key sdk.StoreKey, ibcKeeper *mock.Keeper, stakingKeeper staking.Keeper, supplyKeeper supply.Keeper) StakingMintKeeper {
	keeper := StakingMintKeeper{
		cdc:      cdc,
		storeKey: key,

		ibcKeeper:     *ibcKeeper,
		stakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
	}
	ibcKeeper.SetOnReceivePacket(keeper.OnReceivePacket)
	return keeper
}

func (keeper StakingMintKeeper) OnReceivePacket(ctx sdk.Context, packet []byte) error {
	ibcDelegate := stakingibc.PacketIBCDelegate{}
	err := keeper.cdc.UnmarshalBinaryLengthPrefixed(packet, &ibcDelegate)
	if err != nil {
		return err
	}

	// bondDenom := keeper.stakingKeeper.BondDenom(ctx)

	// TODO: separting for each chain id
	recipientModuleName := "staking-mint"
	coins := sdk.NewCoins(sdk.NewCoin("buatom", ibcDelegate.Amount.Amount))
	err = keeper.supplyKeeper.MintCoins(ctx, recipientModuleName, coins)
	if err != nil {
		return err
	}

	err = keeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, recipientModuleName, sdk.AccAddress(ibcDelegate.DstAddress), coins)
	if err != nil {
		return err
	}

	return nil
}
