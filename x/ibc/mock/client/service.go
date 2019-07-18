package client

import (
	"github.com/tendermint/go-amino"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client/context"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/tendermint/tendermint/rpc/client"
)

type RelayerService struct {
	cmn.BaseService
	cdc *amino.Codec

	watch string

	passphrase string
	txBldr     authtxb.TxBuilder
	cliCtx     context.CLIContext
	client *client.HTTP
}

func NewRelayerService(cdc *amino.Codec, watch string) *RelayerService {
	rs := &RelayerService{
		cdc:   cdc,
		watch: watch,
	}
	rs.BaseService = *cmn.NewBaseService(nil, "RelayerService", rs)
	return rs
}

func (rs *RelayerService) OnStart() error {
	err := rs.init()
	if err != nil {
		return err
	}

	go rs.txRoutine()

	return nil
}
