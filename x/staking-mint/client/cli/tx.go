package cli

import (
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	stakingmint "github.com/cosmos/cosmos-sdk/x/staking-mint"

	"github.com/spf13/cobra"
)

// GetCmdDelegate implements the delegate command.
func GetCmdUndelegate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "undelegate",
		Args:  cobra.ExactArgs(0),
		Short: "undelegate liquid tokens to a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			delAddr := cliCtx.GetFromAddress()

			msg := stakingmint.NewMsgLiquidUndelegate(delAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}
