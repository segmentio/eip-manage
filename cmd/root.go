package cmd

import "github.com/spf13/cobra"

var dryRun bool

var RootCmd = &cobra.Command{
	Use:   "eip-manage",
	Short: "Manage AWS Elastic IP within an EC2 instance",
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false,
		"Dry run")
}
