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
	ibcKeeper.AddOnReceivePacket(keeper.OnReceivePacket)
	return keeper
}

func (keeper StakingMintKeeper) OnReceivePacket(ctx sdk.Context, packet []byte) error {
	store := ctx.KVStore(keeper.storeKey)

	var iPacket mock.Packet
	err := keeper.cdc.UnmarshalBinaryLengthPrefixed(packet, &iPacket)
	if err != nil {
		return err
	}

	ibcDelegate := stakingibc.PacketIBCDelegated{}
	switch iPacket := iPacket.(type) {
	case stakingibc.PacketIBCDelegated:
		ibcDelegate = iPacket
	default:
		return nil
	}

	acc := sdk.AccAddress(ibcDelegate.DstAddress)
	if store.Has(acc.Bytes()) {
		// Due to the lack of development time, limit delegators to delegate only once at a time.
		return sdk.ErrInternal("Can't mint delegatation twice with an account")
	}

	bz, err := keeper.cdc.MarshalBinaryBare(liquidDelegateInfo{
		Delegator:   ibcDelegate.From,
		Validator:   ibcDelegate.Validator,
		DestAddress: ibcDelegate.DstAddress,
		Amount:      ibcDelegate.Amount,
	})
	if err != nil {
		return err
	}
	store.Set(acc.Bytes(), bz)

	// bondDenom := keeper.stakingKeeper.BondDenom(ctx)

	// TODO: separting for each chain id
	recipientModuleName := "staking-mint"
	ratio, err := sdk.NewDecFromStr("0.9")
	if err != nil {
		return err
	}
	amount := sdk.NewDecFromInt(ibcDelegate.Amount.Amount).Mul(ratio).RoundInt()
	coins := sdk.NewCoins(sdk.NewCoin("buatom", amount))
	err = keeper.supplyKeeper.MintCoins(ctx, recipientModuleName, coins)
	if err != nil {
		return err
	}

	err = keeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, recipientModuleName, acc, coins)
	if err != nil {
		return err
	}

	return nil
}

func (keeper StakingMintKeeper) UndelegateLiquid(ctx sdk.Context, sender sdk.AccAddress) (sdk.Tags, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)

	if !store.Has(sender.Bytes()) {
		return sdk.Tags{}, sdk.ErrInternal("Account doesn't have liquid delegation info")
	}

	info := liquidDelegateInfo{}
	err := keeper.cdc.UnmarshalBinaryBare(store.Get(sender.Bytes()), &info)
	if err != nil {
		return sdk.Tags{}, sdk.ErrInternal(err.Error())
	}

	ratio, err := sdk.NewDecFromStr("0.9")
	if err != nil {
		return sdk.Tags{}, sdk.ErrInternal(err.Error())
	}
	amount := sdk.NewDecFromInt(info.Amount.Amount).Mul(ratio).RoundInt()

	sdkErr := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, "staking-burn", sdk.NewCoins(sdk.NewCoin("uatom", amount)))
	if sdkErr != nil {
		return sdk.Tags{}, sdkErr
	}
	sdkErr = keeper.supplyKeeper.BurnCoins(ctx, "staking-burn", sdk.NewCoins(sdk.NewCoin("uatom", amount)))
	if sdkErr != nil {
		return sdk.Tags{}, sdkErr
	}

	bz, err := keeper.cdc.MarshalBinaryLengthPrefixed(stakingibc.PacketIBCUndelegate{
		Delegator: info.Delegator,
		Validator: info.Validator,
		Amount:    info.Amount,
	})
	if err != nil {
		return sdk.Tags{}, sdk.ErrInternal(err.Error())
	}

	tags, sdkErr := keeper.ibcKeeper.SendPacket(ctx, "TODO", mock.MockPacket{
		Data: bz,
	})
	if sdkErr != nil {
		return sdk.Tags{}, sdkErr
	}

	store.Delete(sender.Bytes())

	return tags, nil
}
