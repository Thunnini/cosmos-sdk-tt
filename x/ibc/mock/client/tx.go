package client

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

)

func (rs *RelayerService) init() error {
	rs.txBldr = authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(rs.cdc))
	rs.cliCtx = context.NewCLIContext().
		WithCodec(rs.cdc).
		WithAccountDecoder(rs.cdc)

	fromName := rs.cliCtx.GetFromName()
	_passphrase, err := keys.GetPassphrase(fromName)
	if err != nil {
		return err
	}
	rs.passphrase = _passphrase

	return nil
}

func (rs *RelayerService) txRoutine() {
	// httpClient := client.NewHTTP(rs.cliCtx.NodeURI, "/websocket")

	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					rs.Logger.Error("Unknown error", r)
				}

				time.Sleep(1 * time.Second)
			}()

			rs.Logger.Info("!!!!")
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
