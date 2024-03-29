package cmd

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/segmentio/eip-manage/lib"
	"github.com/spf13/cobra"
)

var inNetwork, outNetwork string
var instanceId string

var associateCmd = &cobra.Command{
	Use:   "associate",
	Short: "Associate an Elastic IP to an EC2 instance",
	RunE:  associate,
}

func init() {
	RootCmd.AddCommand(associateCmd)
	associateCmd.PersistentFlags().StringVarP(&inNetwork, "in-network", "i",
		"", "if set, look for addresses within the network")
	associateCmd.PersistentFlags().StringVarP(&outNetwork, "out-network", "o",
		"", "if set, look for addresses outside the network")
	associateCmd.Flags().StringVar(&instanceId, "instance-id", "",
		"specify the targeted instance id. If not set, use the metadata")
}

func associate(cmd *cobra.Command, args []string) error {
	var ip *ec2.Address
	var err error

	client := lib.NewClient()

	if err := client.SetInstanceId(instanceId); err != nil {
		return err
	}

	if inNetwork != "" {
		ip, err = client.GetIpInNetwork(inNetwork)
	} else if outNetwork != "" {
		ip, err = client.GetIpOutNetwork(inNetwork)
	} else {
		ip, err = client.GetAvailableIP()
	}

	if err != nil {
		return err
	}

	if err := client.AssociateIp(*ip.AllocationId); err != nil {
		return err
	}

	return nil
}
