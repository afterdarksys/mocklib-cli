package main

import (
	"flag"
	"fmt"
	"os"
)

func iamPrintUsage() {
	fmt.Print(`Usage: mocklib-cli iam <action> [flags]

Actions:
  create-user         Create an IAM user
  list-users          List all IAM users
  delete-user         Delete an IAM user
  create-access-key   Create an access key for a user
  create-role         Create an IAM role
  create-policy       Create an IAM managed policy
  attach-user-policy  Attach a managed policy to a user

Flags:
  create-user:
    --name          Username (required)

  delete-user:
    --name          Username (required)

  create-access-key:
    --user          Username (required)

  create-role:
    --name          Role name (required)
    --policy        Trust-policy JSON document (required)

  create-policy:
    --name          Policy name (required)
    --document      Policy JSON document (required)

  attach-user-policy:
    --user          Username (required)
    --policy-arn    Policy ARN to attach (required)
`)
}

func runIAM(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		iamPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create-user":
		fs := flag.NewFlagSet("create-user", flag.ExitOnError)
		name := fs.String("name", "", "Username")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action":   "CreateUser",
			"UserName": *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-users":
		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action": "ListUsers",
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-user":
		fs := flag.NewFlagSet("delete-user", flag.ExitOnError)
		name := fs.String("name", "", "Username")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action":   "DeleteUser",
			"UserName": *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-access-key":
		fs := flag.NewFlagSet("create-access-key", flag.ExitOnError)
		user := fs.String("user", "", "Username")
		fs.Parse(rest)
		requireArg("user", *user)

		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action":   "CreateAccessKey",
			"UserName": *user,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-role":
		fs := flag.NewFlagSet("create-role", flag.ExitOnError)
		name := fs.String("name", "", "Role name")
		policy := fs.String("policy", "", "Trust-policy JSON document")
		fs.Parse(rest)
		requireArg("name", *name)
		requireArg("policy", *policy)

		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action":                  "CreateRole",
			"RoleName":                *name,
			"AssumeRolePolicyDocument": *policy,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-policy":
		fs := flag.NewFlagSet("create-policy", flag.ExitOnError)
		name := fs.String("name", "", "Policy name")
		document := fs.String("document", "", "Policy JSON document")
		fs.Parse(rest)
		requireArg("name", *name)
		requireArg("document", *document)

		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action":         "CreatePolicy",
			"PolicyName":     *name,
			"PolicyDocument": *document,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "attach-user-policy":
		fs := flag.NewFlagSet("attach-user-policy", flag.ExitOnError)
		user := fs.String("user", "", "Username")
		policyARN := fs.String("policy-arn", "", "Policy ARN to attach")
		fs.Parse(rest)
		requireArg("user", *user)
		requireArg("policy-arn", *policyARN)

		resp, err := makeFormRequest("/iam/", map[string]string{
			"Action":    "AttachUserPolicy",
			"UserName":  *user,
			"PolicyArn": *policyARN,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown iam action %q\n\n", action)
		iamPrintUsage()
		os.Exit(1)
	}
}
