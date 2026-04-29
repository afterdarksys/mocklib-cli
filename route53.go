package main

import (
	"flag"
	"fmt"
	"os"
)

func route53PrintUsage() {
	fmt.Print(`Usage: mocklib-cli route53 <action> [flags]

Actions:
  create-zone      Create a hosted zone
  list-zones       List all hosted zones
  change-records   Create/update/delete a DNS record set

Flags:
  create-zone:
    --name    Domain name for the zone (required)

  change-records:
    --zone    Hosted zone ID (required)
    --action  Change action: CREATE | UPSERT | DELETE (required)
    --name    Record name (required)
    --type    Record type, e.g. A, CNAME, MX (required)
    --value   Record value / target (required)
    --ttl     TTL in seconds (default: 300)
`)
}

func runRoute53(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		route53PrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create-zone":
		fs := flag.NewFlagSet("create-zone", flag.ExitOnError)
		name := fs.String("name", "", "Domain name for the hosted zone")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/route53/", map[string]string{
			"Action": "CreateHostedZone",
			"Name":   *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-zones":
		resp, err := makeFormRequest("/route53/", map[string]string{
			"Action": "ListHostedZones",
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "change-records":
		fs := flag.NewFlagSet("change-records", flag.ExitOnError)
		zone := fs.String("zone", "", "Hosted zone ID")
		changeAction := fs.String("action", "", "Change action: CREATE | UPSERT | DELETE")
		name := fs.String("name", "", "Record name")
		recordType := fs.String("type", "", "Record type (A, CNAME, MX, …)")
		value := fs.String("value", "", "Record value / target")
		ttl := fs.String("ttl", "300", "TTL in seconds")
		fs.Parse(rest)
		requireArg("zone", *zone)
		requireArg("action", *changeAction)
		requireArg("name", *name)
		requireArg("type", *recordType)
		requireArg("value", *value)

		resp, err := makeFormRequest("/route53/", map[string]string{
			"Action":       "ChangeResourceRecordSets",
			"HostedZoneId": *zone,
			"ChangeAction": *changeAction,
			"Name":         *name,
			"Type":         *recordType,
			"Value":        *value,
			"TTL":          *ttl,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown route53 action %q\n\n", action)
		route53PrintUsage()
		os.Exit(1)
	}
}
