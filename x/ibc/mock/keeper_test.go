package mock

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPacket(t *testing.T) {
	input := setupInput()
	chain1 := input.chain1
	chain2 := input.chain2

	_, err := chain1.ik.SendPacket(chain1.ctx, chain2.ctx.ChainID(), MockPacket{
		Data: []byte{0},
	})
	require.Nil(t, err)

	err = chain2.ik.ReceivePacket(chain2.ctx, chain1.ctx.ChainID(), MockPacket{
		Data: []byte{0},
	})
	require.Nil(t, err)
}

/*
func TestOpenChannel(t *testing.T) {
	input := setupInput()
	chain1 := input.chain1
	chain2 := input.chain2

	err := chain1.ik.CreateClient(chain1.ctx, "test")
	require.Nil(t, err, "fail to create client")

	err = chain1.ik.CreateClient(chain1.ctx, "test")
	require.NotNil(t, err, "can't create duplicated client")
	require.Contains(t, err.Error(), "Create client on already existing id")

	_, err = chain1.ik.ConnOpenInit(chain1.ctx, "connection", ibc_connection.Connection{
		Client:       "test",
		Counterparty: "counterparty",
		Path:         MockPath{},
	}, "counterpartyclient")
	require.Nil(t, err, "fail to open init")

	err = chain2.ik.CreateClient(chain2.ctx, "test")
	require.Nil(t, err, "fail to create client on chain 2")

	_, err = chain2.ik.ConnOpenTry(chain2.ctx, "connection", ibc_connection.Connection{
		Client:       "test",
		Counterparty: "counterparty",
		Path:         MockPath{},
	}, "counterpartyclient")
	fmt.Println(err.Error())
	require.Nil(t, err, "fail to try open")
}
*/

type input struct {
	chain1 *testChain
	chain2 *testChain
}

type testChain struct {
	cdc *codec.Codec
	ctx sdk.Context
	ik  Keeper
}

func setupInput() input {
	return input{
		chain1: setupChain("test-1"),
		chain2: setupChain("test-2"),
	}
}

func setupChain(chainId string) *testChain {
	db := dbm.NewMemDB()

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)

	ibcKey := sdk.NewKVStoreKey("IBCKey")
	keyParams := sdk.NewKVStoreKey("subspace")
	tkeyParams := sdk.NewTransientStoreKey("transient_subspace")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ik := NewKeeper(cdc, ibcKey, "conn")

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainId}, false, log.NewNopLogger())

	return &testChain{cdc: cdc, ctx: ctx, ik: ik}
}
