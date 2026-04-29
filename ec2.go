package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func ec2PrintUsage() {
	fmt.Print(`Usage: mocklib-cli ec2 <action> [flags]

Actions:
  run-instances        Launch EC2 instances
  describe-instances   List all instances
  start                Start a stopped instance
  stop                 Stop a running instance
  terminate            Terminate an instance

Flags:
  run-instances:
    --type   Instance type (default: t2.micro)
    --ami    AMI ID (default: ami-12345678)
    --count  Number of instances (default: 1)
  start / stop / terminate:
    --id     Instance ID (required)
`)
}

func runEC2(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		ec2PrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "run-instances":
		fs := flag.NewFlagSet("run-instances", flag.ExitOnError)
		instanceType := fs.String("type", "t2.micro", "Instance type")
		ami := fs.String("ami", "ami-12345678", "AMI ID")
		count := fs.Int("count", 1, "Number of instances")
		fs.Parse(rest)

		resp, err := makeFormRequest("/ec2/", map[string]string{
			"Action":       "RunInstances",
			"InstanceType": *instanceType,
			"ImageId":      *ami,
			"MinCount":     strconv.Itoa(*count),
			"MaxCount":     strconv.Itoa(*count),
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "describe-instances":
		resp, err := makeFormRequest("/ec2/", map[string]string{
			"Action": "DescribeInstances",
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "start":
		fs := flag.NewFlagSet("start", flag.ExitOnError)
		id := fs.String("id", "", "Instance ID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := makeFormRequest("/ec2/", map[string]string{
			"Action":      "StartInstances",
			"InstanceId":  *id,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "stop":
		fs := flag.NewFlagSet("stop", flag.ExitOnError)
		id := fs.String("id", "", "Instance ID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := makeFormRequest("/ec2/", map[string]string{
			"Action":      "StopInstances",
			"InstanceId":  *id,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "terminate":
		fs := flag.NewFlagSet("terminate", flag.ExitOnError)
		id := fs.String("id", "", "Instance ID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := makeFormRequest("/ec2/", map[string]string{
			"Action":      "TerminateInstances",
			"InstanceId":  *id,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown ec2 action %q\n\n", action)
		ec2PrintUsage()
		os.Exit(1)
	}
}
