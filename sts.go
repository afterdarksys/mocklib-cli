package main

import (
	"flag"
	"fmt"
	"os"
)

func stsPrintUsage() {
	fmt.Print(`Usage: mocklib-cli sts <action> [flags]

Actions:
  get-caller-identity          Return details about the current identity
  assume-role                  Assume an IAM role

Flags for assume-role:
  --role-arn      ARN of the role to assume (required)
  --session-name  Name for the assumed-role session (required)
`)
}

func runSTS(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		stsPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "get-caller-identity":
		resp, err := makeFormRequest("/sts/", map[string]string{
			"Action": "GetCallerIdentity",
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "assume-role":
		fs := flag.NewFlagSet("assume-role", flag.ExitOnError)
		roleARN := fs.String("role-arn", "", "ARN of the role to assume")
		sessionName := fs.String("session-name", "", "Name for the assumed-role session")
		fs.Parse(rest)
		requireArg("role-arn", *roleARN)
		requireArg("session-name", *sessionName)

		resp, err := makeFormRequest("/sts/", map[string]string{
			"Action":          "AssumeRole",
			"RoleArn":         *roleARN,
			"RoleSessionName": *sessionName,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown sts action %q\n\n", action)
		stsPrintUsage()
		os.Exit(1)
	}
}
