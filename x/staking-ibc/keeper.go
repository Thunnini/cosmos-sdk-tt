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

func NewStakingIBCKeeper(cdc *codec.Codec, key sdk.StoreKey, ibcKeeper mock.Keeper, stakingKeeper staking.Keeper, supplyKeeper supply.Keeper) StakingIBCKeeper {
	return StakingIBCKeeper{
		cdc:      cdc,
		storeKey: key,

		ibcKeeper:     ibcKeeper,
		stakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
	}
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

	bz, gerr = keeper.cdc.MarshalBinaryBare(PacketIBCDelegate{
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

	/*shares, err := keeper.stakingKeeper.ValidateUnbondAmount(
		ctx, moduleAddress, validatorAddress, amount.Amount,
	)*/

	// Due to lack of development time, delegate just can undelegate all right now.
	_, err := keeper.stakingKeeper.Undelegate(ctx, moduleAddress, validatorAddress, sdk.NewDec(1))
	if err != nil {
		return err
	}

	return nil
}
