package main

import (
	"flag"
	"fmt"
	"os"
)

func sqsPrintUsage() {
	fmt.Print(`Usage: mocklib-cli sqs <action> [flags]

Actions:
  create-queue     Create an SQS queue
  send             Send a message to a queue
  receive          Receive messages from a queue
  delete-message   Delete a message from a queue
  delete-queue     Delete a queue
  purge            Purge all messages from a queue

Flags:
  create-queue:
    --name      Queue name (required)

  send:
    --queue     Queue URL (required)
    --message   Message body (required)

  receive:
    --queue     Queue URL (required)
    --max       Maximum number of messages to receive (default: 1)

  delete-message:
    --queue     Queue URL (required)
    --receipt   Receipt handle of the message (required)

  delete-queue:
    --name      Queue name or URL (required)

  purge:
    --queue     Queue URL (required)
`)
}

func runSQS(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		sqsPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create-queue":
		fs := flag.NewFlagSet("create-queue", flag.ExitOnError)
		name := fs.String("name", "", "Queue name")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/aws/sqs", map[string]string{
			"Action":    "CreateQueue",
			"QueueName": *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "send":
		fs := flag.NewFlagSet("send", flag.ExitOnError)
		queue := fs.String("queue", "", "Queue URL")
		message := fs.String("message", "", "Message body")
		fs.Parse(rest)
		requireArg("queue", *queue)
		requireArg("message", *message)

		resp, err := makeFormRequest("/aws/sqs", map[string]string{
			"Action":      "SendMessage",
			"QueueUrl":    *queue,
			"MessageBody": *message,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "receive":
		fs := flag.NewFlagSet("receive", flag.ExitOnError)
		queue := fs.String("queue", "", "Queue URL")
		max := fs.String("max", "1", "Maximum number of messages to receive")
		fs.Parse(rest)
		requireArg("queue", *queue)

		resp, err := makeFormRequest("/aws/sqs", map[string]string{
			"Action":              "ReceiveMessage",
			"QueueUrl":            *queue,
			"MaxNumberOfMessages": *max,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-message":
		fs := flag.NewFlagSet("delete-message", flag.ExitOnError)
		queue := fs.String("queue", "", "Queue URL")
		receipt := fs.String("receipt", "", "Receipt handle")
		fs.Parse(rest)
		requireArg("queue", *queue)
		requireArg("receipt", *receipt)

		resp, err := makeFormRequest("/aws/sqs", map[string]string{
			"Action":        "DeleteMessage",
			"QueueUrl":      *queue,
			"ReceiptHandle": *receipt,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-queue":
		fs := flag.NewFlagSet("delete-queue", flag.ExitOnError)
		name := fs.String("name", "", "Queue name or URL")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/aws/sqs", map[string]string{
			"Action":   "DeleteQueue",
			"QueueUrl": *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "purge":
		fs := flag.NewFlagSet("purge", flag.ExitOnError)
		queue := fs.String("queue", "", "Queue URL")
		fs.Parse(rest)
		requireArg("queue", *queue)

		resp, err := makeFormRequest("/aws/sqs", map[string]string{
			"Action":   "PurgeQueue",
			"QueueUrl": *queue,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown sqs action %q\n\n", action)
		sqsPrintUsage()
		os.Exit(1)
	}
}
