package main

import (
	"flag"
	"fmt"
	"os"
)

func snsPrintUsage() {
	fmt.Print(`Usage: mocklib-cli sns <action> [flags]

Actions:
  create-topic   Create an SNS topic
  list-topics    List all SNS topics
  publish        Publish a message to a topic
  subscribe      Subscribe an endpoint to a topic

Flags:
  create-topic:
    --name        Topic name (required)

  publish:
    --topic       Topic ARN (required)
    --message     Message body (required)

  subscribe:
    --topic       Topic ARN (required)
    --protocol    Subscription protocol, e.g. email, sqs, http (required)
    --endpoint    Subscription endpoint (required)
`)
}

func runSNS(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		snsPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create-topic":
		fs := flag.NewFlagSet("create-topic", flag.ExitOnError)
		name := fs.String("name", "", "Topic name")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/sns/", map[string]string{
			"Action": "CreateTopic",
			"Name":   *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-topics":
		resp, err := makeFormRequest("/sns/", map[string]string{
			"Action": "ListTopics",
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "publish":
		fs := flag.NewFlagSet("publish", flag.ExitOnError)
		topic := fs.String("topic", "", "Topic ARN")
		message := fs.String("message", "", "Message body")
		fs.Parse(rest)
		requireArg("topic", *topic)
		requireArg("message", *message)

		resp, err := makeFormRequest("/sns/", map[string]string{
			"Action":   "Publish",
			"TopicArn": *topic,
			"Message":  *message,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "subscribe":
		fs := flag.NewFlagSet("subscribe", flag.ExitOnError)
		topic := fs.String("topic", "", "Topic ARN")
		protocol := fs.String("protocol", "", "Subscription protocol (email, sqs, http, …)")
		endpoint := fs.String("endpoint", "", "Subscription endpoint")
		fs.Parse(rest)
		requireArg("topic", *topic)
		requireArg("protocol", *protocol)
		requireArg("endpoint", *endpoint)

		resp, err := makeFormRequest("/sns/", map[string]string{
			"Action":   "Subscribe",
			"TopicArn": *topic,
			"Protocol": *protocol,
			"Endpoint": *endpoint,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown sns action %q\n\n", action)
		snsPrintUsage()
		os.Exit(1)
	}
}
