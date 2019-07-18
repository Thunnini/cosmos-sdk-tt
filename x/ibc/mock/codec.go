package mock

import (
	"github.com/cosmos/cosmos-sdk/codec"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	ibc_channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/merkle"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*client.ConsensusState)(nil), nil)
	cdc.RegisterInterface((*commitment.Path)(nil), nil)
	cdc.RegisterInterface((*commitment.Root)(nil), nil)
	cdc.RegisterInterface((*commitment.Proof)(nil), nil)
	cdc.RegisterInterface((*ibc_channel.Packet)(nil), nil)

	cdc.RegisterConcrete(MsgCreateClient{}, "cosmos-sdk/MsgCreateClient", nil)
	cdc.RegisterConcrete(MockConsensusState{}, "mockconsensusstate", nil)
	cdc.RegisterConcrete(MockPath{}, "mockpath", nil)
	cdc.RegisterConcrete(MockRoot{}, "mockroot", nil)
	cdc.RegisterConcrete(MockProof{}, "mockproof", nil)
	cdc.RegisterConcrete(merkle.Path{}, "merklepath", nil)
	cdc.RegisterConcrete(MockPacket{}, "mockpacket", nil)
}

// generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
