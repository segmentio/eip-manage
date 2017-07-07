package lib

import (
	"errors"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type Client struct {
	Svc        ec2iface.EC2API
	Metadata   *ec2metadata.EC2Metadata
	InstanceId string
}

func NewClient() *Client {
	session := session.New()
	return &Client{
		Svc:      ec2.New(session),
		Metadata: ec2metadata.New(session),
	}
}

func (c *Client) SetInstanceId(id string) error {
	var err error
	if id == "" {
		id, err = c.Metadata.GetMetadata("instance-id")
		if err != nil {
			return err
		}
	}
	c.InstanceId = id
	return nil
}

func (c *Client) AssociateIp(eipalloc string) error {
	input := &ec2.AssociateAddressInput{
		AllocationId: aws.String(eipalloc),
		InstanceId:   aws.String(c.InstanceId),
	}

	//TODO: Add wait + timeout here
	_, err := c.Svc.AssociateAddress(input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetIpInNetwork(network string) (*ec2.Address, error) {
	_, n, _ := net.ParseCIDR(network)

	filter := func(ips []*ec2.Address) (result []*ec2.Address) {
		for _, ip := range ips {
			if ip.AssociationId == nil && n.Contains(net.ParseIP(*ip.PublicIp)) {
				result = append(result, ip)
				break
			}
		}
		return
	}

	eips, err := c.getEIPs(filter)
	if err != nil {
		return nil, err
	}

	if len(eips) < 1 {
		return nil, errors.New("no ip available")
	}

	return eips[0], nil
}

func (c *Client) GetIpOutNetwork(network string) (*ec2.Address, error) {
	_, n, _ := net.ParseCIDR(network)

	filter := func(ips []*ec2.Address) (result []*ec2.Address) {
		for _, ip := range ips {
			if ip.AssociationId == nil && n.Contains(net.ParseIP(*ip.PublicIp)) == false {
				result = append(result, ip)
				break
			}
		}
		return
	}

	eips, err := c.getEIPs(filter)
	if err != nil {
		return nil, err
	}

	if len(eips) < 1 {
		return nil, errors.New("no ip available")
	}

	return eips[0], nil
}

func (c *Client) GetAvailableIP() (*ec2.Address, error) {
	filter := func(ips []*ec2.Address) (result []*ec2.Address) {
		for _, ip := range ips {
			if ip.AssociationId == nil {
				result = append(result, ip)
				break
			}
		}
		return
	}

	eips, err := c.getEIPs(filter)
	if err != nil {
		return nil, err
	}

	if len(eips) < 1 {
		return nil, errors.New("no ip available")
	}

	return eips[0], nil
}

func (c *Client) GetAvailableIPs() ([]*ec2.Address, error) {
	filter := func(ips []*ec2.Address) (result []*ec2.Address) {
		for _, ip := range ips {
			if ip.AssociationId == nil {
				result = append(result, ip)
			}
		}
		return
	}

	return c.getEIPs(filter)
}

func (c *Client) getEIPs(filter func([]*ec2.Address) []*ec2.Address) ([]*ec2.Address, error) {
	var addresses []*ec2.Address
	input := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			ec2Filter("domain", "vpc"),
		},
	}

	result, err := c.Svc.DescribeAddresses(input)
	if err != nil {
		return nil, err
	}

	if len(result.Addresses) > 0 {
		addresses = filter(result.Addresses)
	}

	return addresses, nil
}

func ec2Filter(name string, values ...string) *ec2.Filter {
	return &ec2.Filter{
		Name:   &name,
		Values: aws.StringSlice(values),
	}
}

func (c *Client) getInstanceId() (string, error) {
	return c.Metadata.GetMetadata("instance-id")
}
