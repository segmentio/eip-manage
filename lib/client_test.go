package lib

import "github.com/aws/aws-sdk-go/service/ec2/ec2iface"

type mockEC2Client struct {
	ec2iface.EC2API
}
