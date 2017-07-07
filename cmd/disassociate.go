package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var disassociateCmd = &cobra.Command{
	Use:   "disassociate",
	Short: "Disassociate an Elastic IP from an EC2 instance",
	RunE:  disassociate,
}

func init() {
	RootCmd.AddCommand(disassociateCmd)
}

func disassociate(cmd *cobra.Command, args []string) error {
	fmt.Println("Not yet implemented.")
	return nil
}
