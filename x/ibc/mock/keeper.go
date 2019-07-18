package mock

import (
	"encoding/base64"
	"math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/state"
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	ibc_connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	ibc_channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

// TODO: add router
type Keeper struct {
	cdc              *codec.Codec
	key              sdk.StoreKey
	connId           string
	onPacketReceives []func(sdk.Context, []byte) error

	clientMan client.Manager

	connMan        ibc_connection.Manager
	connHandshaker ibc_connection.Handshaker

	chMan        ibc_channel.Manager
	chHandshaker ibc_channel.Handshaker
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, connId string) Keeper {
	base := state.NewBase(cdc, key, []byte{})
	clientMan := client.NewManager(base, base)
	connMan := ibc_connection.NewManager(base, client.NewManager(base, base))
	chMan := ibc_channel.NewManager(base, connMan)

	return Keeper{
		cdc:              cdc,
		key:              key,
		connId:           connId,
		onPacketReceives: make([]func(sdk.Context, []byte) error, 0),

		clientMan: clientMan,

		connMan:        connMan,
		connHandshaker: ibc_connection.NewHandshaker(connMan),

		chMan:        chMan,
		chHandshaker: ibc_channel.NewHandshaker(chMan),
	}
}

func (keeper *Keeper) AddOnReceivePacket(onPacketReceive func(sdk.Context, []byte) error) {
	keeper.onPacketReceives = append(keeper.onPacketReceives, onPacketReceive)
}

func (keeper Keeper) SendPacket(ctx sdk.Context, counterChainId string, packet MockPacket) (sdk.Tags, sdk.Error) {
	cobj := ibc_channel.NewCounterObject(keeper.cdc, keeper.key, counterChainId, keeper.connId)
	obj := ibc_channel.NewObject(keeper.cdc, keeper.key, ctx.ChainID(), keeper.connId, cobj)
	err := obj.Send(ctx, packet)
	if err != nil {
		return sdk.Tags{}, sdk.ErrInternal(err.Error())
	}

	return sdk.NewTags("type", "ibc-send", "data", base64.StdEncoding.EncodeToString(packet.Data)), nil
}

func (keeper Keeper) ReceivePacket(ctx sdk.Context, counterChainId string, packet MockPacket) sdk.Error {
	cobj := ibc_channel.NewCounterObject(keeper.cdc, keeper.key, counterChainId, keeper.connId)
	obj := ibc_channel.NewObject(keeper.cdc, keeper.key, ctx.ChainID(), keeper.connId, cobj)
	err := obj.Receive(ctx, MockProof{}, packet)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}

	for _, onPacketReceive := range keeper.onPacketReceives {
		err := onPacketReceive(ctx, packet.Data)
		if err != nil {
			return sdk.ErrInternal(err.Error())
		}
	}

	return nil
}

func (keeper Keeper) CreateClient(ctx sdk.Context, clientId string) sdk.Error {
	_, err := keeper.clientMan.Create(ctx, clientId, MockConsensusState{})
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}

	return nil
}

func (keeper Keeper) ConnOpenInit(ctx sdk.Context, id string, connection ibc_connection.Connection, counterpartyClient string) (ibc_connection.HandshakeObject, sdk.Error) {
	obj, err := keeper.connHandshaker.OpenInit(ctx, id, connection, counterpartyClient, math.MaxUint64)
	if err != nil {
		return ibc_connection.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ConnOpenTry(ctx sdk.Context, id string, connection ibc_connection.Connection, counterpartyClient string) (ibc_connection.HandshakeObject, sdk.Error) {
	obj, err := keeper.connHandshaker.OpenTry(ctx, MockProof{}, MockProof{}, MockProof{}, MockProof{}, id, connection, counterpartyClient, math.MaxUint64, math.MaxUint64)
	if err != nil {
		return ibc_connection.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ConnOpenAck(ctx sdk.Context, id string) (ibc_connection.HandshakeObject, sdk.Error) {
	obj, err := keeper.connHandshaker.OpenAck(ctx, MockProof{}, MockProof{}, MockProof{}, MockProof{}, id, math.MaxUint64, math.MaxUint64)
	if err != nil {
		return ibc_connection.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ConnOpenConfirm(ctx sdk.Context, statep, timeoutp commitment.Proof, id string, timeoutHeight uint64) (ibc_connection.HandshakeObject, sdk.Error) {
	obj, err := keeper.connHandshaker.OpenConfirm(ctx, MockProof{}, MockProof{}, id, math.MaxUint64)
	if err != nil {
		return ibc_connection.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ChanOpenInit(ctx sdk.Context, connid, chanid string, channel ibc_channel.Channel) (ibc_channel.HandshakeObject, sdk.Error) {
	obj, err := keeper.chHandshaker.OpenInit(ctx, connid, chanid, channel, math.MaxUint64)
	if err != nil {
		return ibc_channel.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ChanOpenTry(ctx sdk.Context,
	pchannel, pstate, ptimeout commitment.Proof,
	connid, chanid string, channel ibc_channel.Channel) (ibc_channel.HandshakeObject, sdk.Error) {
	obj, err := keeper.chHandshaker.OpenTry(ctx, pchannel, pstate, ptimeout, connid, chanid, channel, math.MaxUint64, math.MaxUint64)
	if err != nil {
		return ibc_channel.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ChanOpenAck(ctx sdk.Context,
	pchannel, pstate, ptimeout commitment.Proof,
	connid, chanid string) (ibc_channel.HandshakeObject, sdk.Error) {
	obj, err := keeper.chHandshaker.OpenAck(ctx, pchannel, pstate, ptimeout, connid, chanid, math.MaxUint64, math.MaxUint64)
	if err != nil {
		return ibc_channel.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}

func (keeper Keeper) ChanOpenConfirm(ctx sdk.Context,
	pstate, ptimeout commitment.Proof,
	connid, chanid string) (ibc_channel.HandshakeObject, sdk.Error) {
	obj, err := keeper.chHandshaker.OpenConfirm(ctx, pstate, ptimeout, connid, chanid, math.MaxUint64)
	if err != nil {
		return ibc_channel.HandshakeObject{}, sdk.ErrInternal(err.Error())
	}

	return obj, nil
}
