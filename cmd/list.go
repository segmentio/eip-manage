package cmd

import (
	"fmt"

	"github.com/segmentio/eip-manage/lib"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Elastic IPs",
	RunE:  list,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	client := lib.NewClient()

	eips, err := client.GetAvailableIPs()
	if err != nil {
		return err
	}

	for _, eip := range eips {
		fmt.Printf("%s %s\n", *eip.AllocationId, *eip.PublicIp)
	}

	return nil
}
