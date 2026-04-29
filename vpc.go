package main

import (
	"flag"
	"fmt"
	"os"
)

func vpcPrintUsage() {
	fmt.Print(`Usage: mocklib-cli vpc <action> [flags]

Actions:
  create              Create a VPC
  list                List all VPCs
  delete              Delete a VPC
  create-subnet       Create a subnet in a VPC
  create-sg           Create a security group
  authorize-ingress   Add an inbound rule to a security group
  create-igw          Create an internet gateway
  attach-igw          Attach an internet gateway to a VPC
  create-rtb          Create a route table
  create-route        Add a route to a route table
  associate-rtb       Associate a route table with a subnet

Flags:
  create:
    --cidr    CIDR block, e.g. 10.0.0.0/16 (required)

  delete:
    --id      VPC ID (required)

  create-subnet:
    --vpc     VPC ID (required)
    --cidr    Subnet CIDR block (required)
    --az      Availability zone, e.g. us-east-1a (required)

  create-sg:
    --vpc     VPC ID (required)
    --name    Security group name (required)
    --desc    Description (required)

  authorize-ingress:
    --sg        Security group ID (required)
    --protocol  Protocol, e.g. tcp, udp, icmp (required)
    --from      From port (required)
    --to        To port (required)
    --cidr      Source CIDR (required)

  create-igw:        (no flags)

  attach-igw:
    --igw   Internet gateway ID (required)
    --vpc   VPC ID (required)

  create-rtb:
    --vpc   VPC ID (required)

  create-route:
    --rtb   Route table ID (required)
    --cidr  Destination CIDR (required)
    --gw    Gateway ID (required)

  associate-rtb:
    --rtb     Route table ID (required)
    --subnet  Subnet ID (required)
`)
}

// vpcForm sends a form POST to /aws/vpc.
func vpcForm(fields map[string]string) (map[string]interface{}, error) {
	return makeFormRequest("/aws/vpc", fields)
}

func runVPC(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		vpcPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		cidr := fs.String("cidr", "", "CIDR block")
		fs.Parse(rest)
		requireArg("cidr", *cidr)

		resp, err := vpcForm(map[string]string{
			"Action":    "CreateVpc",
			"CidrBlock": *cidr,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list":
		resp, err := vpcForm(map[string]string{"Action": "DescribeVpcs"})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete":
		fs := flag.NewFlagSet("delete", flag.ExitOnError)
		id := fs.String("id", "", "VPC ID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := vpcForm(map[string]string{
			"Action": "DeleteVpc",
			"VpcId":  *id,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-subnet":
		fs := flag.NewFlagSet("create-subnet", flag.ExitOnError)
		vpc := fs.String("vpc", "", "VPC ID")
		cidr := fs.String("cidr", "", "Subnet CIDR block")
		az := fs.String("az", "", "Availability zone")
		fs.Parse(rest)
		requireArg("vpc", *vpc)
		requireArg("cidr", *cidr)
		requireArg("az", *az)

		resp, err := vpcForm(map[string]string{
			"Action":           "CreateSubnet",
			"VpcId":            *vpc,
			"CidrBlock":        *cidr,
			"AvailabilityZone": *az,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-sg":
		fs := flag.NewFlagSet("create-sg", flag.ExitOnError)
		vpc := fs.String("vpc", "", "VPC ID")
		name := fs.String("name", "", "Security group name")
		desc := fs.String("desc", "", "Description")
		fs.Parse(rest)
		requireArg("vpc", *vpc)
		requireArg("name", *name)
		requireArg("desc", *desc)

		resp, err := vpcForm(map[string]string{
			"Action":      "CreateSecurityGroup",
			"VpcId":       *vpc,
			"GroupName":   *name,
			"Description": *desc,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "authorize-ingress":
		fs := flag.NewFlagSet("authorize-ingress", flag.ExitOnError)
		sg := fs.String("sg", "", "Security group ID")
		protocol := fs.String("protocol", "", "Protocol (tcp, udp, icmp)")
		from := fs.String("from", "", "From port")
		to := fs.String("to", "", "To port")
		cidr := fs.String("cidr", "", "Source CIDR")
		fs.Parse(rest)
		requireArg("sg", *sg)
		requireArg("protocol", *protocol)
		requireArg("from", *from)
		requireArg("to", *to)
		requireArg("cidr", *cidr)

		resp, err := vpcForm(map[string]string{
			"Action":     "AuthorizeSecurityGroupIngress",
			"GroupId":    *sg,
			"IpProtocol": *protocol,
			"FromPort":   *from,
			"ToPort":     *to,
			"CidrIp":     *cidr,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-igw":
		resp, err := vpcForm(map[string]string{"Action": "CreateInternetGateway"})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "attach-igw":
		fs := flag.NewFlagSet("attach-igw", flag.ExitOnError)
		igw := fs.String("igw", "", "Internet gateway ID")
		vpc := fs.String("vpc", "", "VPC ID")
		fs.Parse(rest)
		requireArg("igw", *igw)
		requireArg("vpc", *vpc)

		resp, err := vpcForm(map[string]string{
			"Action":            "AttachInternetGateway",
			"InternetGatewayId": *igw,
			"VpcId":             *vpc,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-rtb":
		fs := flag.NewFlagSet("create-rtb", flag.ExitOnError)
		vpc := fs.String("vpc", "", "VPC ID")
		fs.Parse(rest)
		requireArg("vpc", *vpc)

		resp, err := vpcForm(map[string]string{
			"Action": "CreateRouteTable",
			"VpcId":  *vpc,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-route":
		fs := flag.NewFlagSet("create-route", flag.ExitOnError)
		rtb := fs.String("rtb", "", "Route table ID")
		cidr := fs.String("cidr", "", "Destination CIDR")
		gw := fs.String("gw", "", "Gateway ID")
		fs.Parse(rest)
		requireArg("rtb", *rtb)
		requireArg("cidr", *cidr)
		requireArg("gw", *gw)

		resp, err := vpcForm(map[string]string{
			"Action":               "CreateRoute",
			"RouteTableId":         *rtb,
			"DestinationCidrBlock": *cidr,
			"GatewayId":            *gw,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "associate-rtb":
		fs := flag.NewFlagSet("associate-rtb", flag.ExitOnError)
		rtb := fs.String("rtb", "", "Route table ID")
		subnet := fs.String("subnet", "", "Subnet ID")
		fs.Parse(rest)
		requireArg("rtb", *rtb)
		requireArg("subnet", *subnet)

		resp, err := vpcForm(map[string]string{
			"Action":       "AssociateRouteTable",
			"RouteTableId": *rtb,
			"SubnetId":     *subnet,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown vpc action %q\n\n", action)
		vpcPrintUsage()
		os.Exit(1)
	}
}
