package main

import (
	"github.com/spf13/cobra"
	"github.com/teleport-network/teleport/x/aggregate/types"
)

func AddIBCDenomCommand(debug *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-denom [port] [channel] [denom]",
		Short: "Generate ibc denom name",
		Long:  `According to the target channel, port and denom provided by the user, generate the denom name after the ibc cross-chain transfer`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			denom, err := types.IBCDenom(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			cmd.Printf("IBC Denom %s\n for port:%s ,channel:%s ,denom:%s ", denom, args[0], args[1], args[2])
			return nil
		},
	}
	debug.AddCommand(cmd)
	return debug
}
