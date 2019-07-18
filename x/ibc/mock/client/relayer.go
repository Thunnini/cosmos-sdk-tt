package client

import (
	"fmt"
	"os"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/spf13/cobra"
)


var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func GetCmdRelayer(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:  "watch [node]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rs := NewRelayerService(cdc)
			rs.SetLogger(logger.With("module", "relayer"))

			// Stop upon receiving SIGTERM or CTRL-C.
			cmn.TrapSignal(logger, func() {
				if rs.IsRunning() {
					rs.Stop()
				}
			})

			if err := rs.Start(); err != nil {
				return fmt.Errorf("Failed to start node: %v", err)
			}

			// Run forever.
			select {}
		},
	}
}
