package stakingibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type StakingIBCKeeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey

	ibcKeeper     mock.Keeper
	stakingKeeper staking.Keeper
	supplyKeeper  supply.Keeper
}

func NewStakingIBCKeeper(cdc *codec.Codec, key sdk.StoreKey, ibcKeeper *mock.Keeper, stakingKeeper staking.Keeper, supplyKeeper supply.Keeper) StakingIBCKeeper {
	keeper := StakingIBCKeeper{
		cdc:      cdc,
		storeKey: key,

		ibcKeeper:     *ibcKeeper,
		stakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
	}
	ibcKeeper.AddOnReceivePacket(keeper.OnReceivePacket)

	return keeper
}

func (keeper StakingIBCKeeper) Delegate(ctx sdk.Context, from sdk.AccAddress, validatorAddr sdk.ValAddress, amount sdk.Coin, destAddr []byte) (sdk.Tags, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	key := from.Bytes()
	// Due to the lack of development time, limit delegators to delegate only once at a time.
	if store.Has(key) {
		return sdk.Tags{}, sdk.ErrInternal("Can't delegate twice")
	}
	bz, gerr := keeper.cdc.MarshalBinaryBare(delegateInfo{
		Delegator:   from,
		Validator:   validatorAddr,
		DestAddress: destAddr,
	})
	if gerr != nil {
		return sdk.Tags{}, sdk.ErrInternal(gerr.Error())
	}
	store.Set(key, bz)

	bondDenom := keeper.stakingKeeper.BondDenom(ctx)
	if amount.Denom != bondDenom {
		return sdk.Tags{}, sdk.ErrInternal("Invalid denom")
	}
	// TODO: separting for each chain id
	recipientModuleName := "staking-ibc"
	err := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, recipientModuleName, sdk.NewCoins(amount))
	if err != nil {
		return sdk.Tags{}, err
	}

	validator, found := keeper.stakingKeeper.GetValidator(ctx, validatorAddr)
	if !found {
		return sdk.Tags{}, sdk.ErrInternal("Unkwon validator")
	}

	_, err = keeper.stakingKeeper.Delegate(ctx, keeper.supplyKeeper.GetModuleAddress(recipientModuleName), amount.Amount, validator, true)

	bz, gerr = keeper.cdc.MarshalBinaryLengthPrefixed(PacketIBCDelegated{
		From:       from,
		Validator:  validatorAddr,
		Amount:     amount,
		DstAddress: destAddr,
	})
	if gerr != nil {
		return sdk.Tags{}, sdk.ErrInternal(gerr.Error())
	}

	tags, err := keeper.ibcKeeper.SendPacket(ctx, "TODO", mock.MockPacket{
		Data: bz,
	})

	if err != nil {
		return sdk.Tags{}, err
	}

	return tags, nil
}

func (keeper StakingIBCKeeper) Undelegate(ctx sdk.Context, recipient sdk.AccAddress, validatorAddress sdk.ValAddress /*, amount sdk.Coin*/) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	key := recipient.Bytes()
	store.Delete(key)

	/*bondDenom := keeper.stakingKeeper.BondDenom(ctx)
	if amount.Denom != bondDenom {
		return sdk.ErrInternal("Invalid denom")
	}*/

	// TODO: separting for each chain id
	recipientModuleName := "staking-ibc"
	moduleAddress := keeper.supplyKeeper.GetModuleAddress(recipientModuleName)

	del, found := keeper.stakingKeeper.GetDelegation(ctx, moduleAddress, validatorAddress)
	if !found {
		return sdk.ErrInternal("Delegation doesn't exist")
	}

	// Due to lack of development time, delegate just can undelegate all shares right now.
	_, err := keeper.stakingKeeper.Undelegate(ctx, moduleAddress, validatorAddress, del.Shares)
	if err != nil {
		return err
	}

	return nil
}

func (keeper StakingIBCKeeper) OnReceivePacket(ctx sdk.Context, packet []byte) error {
	store := ctx.KVStore(keeper.storeKey)

	var iPacket mock.Packet
	err := keeper.cdc.UnmarshalBinaryLengthPrefixed(packet, &iPacket)
	if err != nil {
		return err
	}

	ibcUndelegate := PacketIBCUndelegate{}
	switch iPacket := iPacket.(type) {
	case PacketIBCUndelegate:
		ibcUndelegate = iPacket
	default:
		return nil
	}

	acc := sdk.AccAddress(ibcUndelegate.Delegator)
	if store.Has(acc.Bytes()) {
		// Due to the lack of development time, limit delegators to delegate only once at a time.
		return sdk.ErrInternal("Can't mint delegatation twice with an account")
	}

	return keeper.Undelegate(ctx, ibcUndelegate.Delegator, ibcUndelegate.Validator)
}
