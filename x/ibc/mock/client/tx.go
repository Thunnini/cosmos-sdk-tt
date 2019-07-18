package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	clictx "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/tendermint/tendermint/rpc/client"

	"github.com/cosmos/cosmos-sdk/x/ibc/mock"
)

func (rs *RelayerService) init() error {
	rs.txBldr = authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(rs.cdc))
	rs.cliCtx = clictx.NewCLIContext().
		WithCodec(rs.cdc).
		WithAccountDecoder(rs.cdc)

	fromName := rs.cliCtx.GetFromName()
	_passphrase, err := keys.GetPassphrase(fromName)
	if err != nil {
		return err
	}
	rs.passphrase = _passphrase

	rs.client = client.NewHTTP(rs.watch, "/websocket")

	return rs.client.OnStart()
}

func (rs *RelayerService) txRoutine() {
	for {
		func() {
			defer func() {
				/*if r := recover(); r != nil {
					rs.Logger.Error("Unknown error", r)
				}*/

				time.Sleep(1 * time.Second)
			}()

			out, err := rs.client.Subscribe(context.Background(), "", "type='ibc-send'")
			if err != nil {
				rs.Logger.Error(err.Error())
			}

			rs.Logger.Info("!!!")

			for true {
				result, ok := <-out
				if !ok {
					rs.Logger.Error("Out channel closed")
					break
				}

				if result.Tags["type"] == "ibc-send" {
					data, err := base64.StdEncoding.DecodeString((result.Tags["data"]))
					if err != nil {
						rs.Logger.Error(err.Error())
						break
					}

					rs.Logger.Info(fmt.Sprintf("Try relay packet %s", result.Tags["data"]))
					relayer := rs.cliCtx.GetFromAddress()
					msg := mock.NewMsgRelay(relayer, data)
					_, err = rs.broadcast([]sdk.Msg{msg})

					if err != nil {
						rs.Logger.Error(err.Error())
						break
					}

					rs.Logger.Info("Succeed to relay")
				}
			}
		}()
	}
}

func (rs *RelayerService) broadcast(msgs []sdk.Msg) (*sdk.TxResponse, error) {
	txBldr, err := utils.PrepareTxBuilder(rs.txBldr, rs.cliCtx)
	if err != nil {
		return nil, err
	}

	fromName := rs.cliCtx.GetFromName()

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, rs.passphrase, msgs)
	if err != nil {
		return nil, err
	}

	// broadcast to a Tendermint node
	res, err := rs.cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	return &res, rs.cliCtx.PrintOutput(res)
}
